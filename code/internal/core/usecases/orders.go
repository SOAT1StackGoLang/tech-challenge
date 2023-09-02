package usecases

import (
	"context"
	"github.com/shopspring/decimal"
	"time"

	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ordersUseCase struct {
	logger     *zap.SugaredLogger
	ordersRepo ports.OrdersRepository
	userUC     ports.UsersUseCase
	prodUC     ports.ProductsUseCase
	paymentsUC ports.PaymentUseCase
}

func NewOrdersUseCase(
	logger *zap.SugaredLogger,
	ordersRepo ports.OrdersRepository,
	userUC ports.UsersUseCase,
	prodUC ports.ProductsUseCase,
	paymentsUC ports.PaymentUseCase,
) ports.OrdersUseCase {
	orderUC := &ordersUseCase{logger: logger, ordersRepo: ordersRepo, userUC: userUC, prodUC: prodUC, paymentsUC: paymentsUC}

	domain.PaymentStatusChannel = make(chan domain.PaymentStatusNotification)
	go orderUC.SubscribeToPaymentStatusUpdates()

	return orderUC
}

func (o *ordersUseCase) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*domain.Order, error) {
	return o.ordersRepo.GetOrderByPaymentID(ctx, paymentID)
}

func (o *ordersUseCase) UpdateOrderStatus(ctx context.Context, userID, orderID uuid.UUID, status domain.OrderStatus) (*domain.Order, error) {
	order, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	return o.ordersRepo.UpdateOrder(ctx, order)
}

func (o *ordersUseCase) Checkout(ctx context.Context, userID, orderID uuid.UUID) (*domain.Order, error) {
	var order *domain.Order

	order, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	order.Status = domain.ORDER_STATUS_WAITING_PAYMENT

	payment, err := o.paymentsUC.CreatePayment(ctx, order)
	if err != nil {
		return nil, err
	}
	order.PaymentID = payment.ID

	defer func() {
		// non blocking step
		order.UpdatedAt = time.Now()
		order, err = o.ordersRepo.UpdateOrder(ctx, order)
		if err != nil {
			o.logger.Errorw(
				"failed updating order status after checkout",
				zap.Error(err),
			)
		}
	}()

	return order, err

}

func (o *ordersUseCase) ListOrders(ctx context.Context, limit, offset int, userID uuid.UUID) (*domain.OrderList, error) {
	if isAdmin(o.logger, o.userUC, ctx, userID) {
		return o.ordersRepo.ListOrders(ctx, limit, offset)
	}
	return o.ordersRepo.ListOrdersByUser(ctx, limit, offset, userID)
}

func (o *ordersUseCase) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*domain.Order, error) {
	order, err := o.ordersRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID && !isAdmin(o.logger, o.userUC, ctx, order.UserID) {
		return nil, helpers.ErrUnauthorized
	}

	return order, nil
}

func (o *ordersUseCase) CreateOrder(ctx context.Context, userID uuid.UUID, products []domain.Product) (*domain.Order, error) {
	var order *domain.Order

	if len(products) == 0 {
		o.logger.Errorw(
			"error at CreateOrder, must have at least one product in it",
			zap.Any("products", products),
			zap.Error(helpers.ErrInvalidInput),
		)
		return nil, helpers.ErrInvalidInput
	}

	for k, p := range products {
		fullProduct, err := o.prodUC.GetProduct(ctx, p.ID)
		if err != nil {
			o.logger.Errorw("CreateOrder failed due to invalid product",
				zap.String("product_id", p.ID.String()),
				zap.Any("requested_products", products),
				zap.Error(err),
			)
			return nil, err
		}
		products[k] = *fullProduct
	}

	order = domain.NewOrder(uuid.New(), userID, time.Now(), products)

	for _, v := range products {
		order.Price = order.Price.Add(v.Price)
	}

	return o.ordersRepo.CreateOrder(ctx, order)
}

func (o *ordersUseCase) InsertProductsIntoOrder(ctx context.Context, userID, orderID uuid.UUID, inProducts []domain.Product) (*domain.Order, error) {
	order, err := o.GetOrder(ctx, userID, orderID)
	// Check ownership
	if err != nil {
		return nil, err
	}

	if len(inProducts) == 0 {
		o.logger.Errorw(
			"error at InsertProductsIntoOrder, must have at least one product in it",
			zap.Any("inProducts", inProducts),
			zap.Error(helpers.ErrInvalidInput),
		)
		return nil, helpers.ErrInvalidInput
	}

	for _, v := range inProducts {
		order.Products = append(order.Products, v)
		order.Price = order.Price.Add(v.Price)
	}

	return o.ordersRepo.UpdateOrder(ctx, order)
}

func (o *ordersUseCase) RemoveProductFromOrder(ctx context.Context, userID, orderID uuid.UUID, outProducts []domain.Product) (*domain.Order, error) {
	order, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	if len(outProducts) == 0 {
		o.logger.Errorw(
			"error at RemoveProductFromOrder, must have at least one product in it",
			zap.Any("outProducts", outProducts),
			zap.Error(helpers.ErrInvalidInput),
		)
		return nil, helpers.ErrInvalidInput
	}

	var removeSet map[uuid.UUID]bool
	for _, p := range outProducts {
		removeSet[p.ID] = true
	}

	var newProdsList []domain.Product
	order.Price = decimal.NewFromInt(0)
	for _, p := range order.Products {
		if _, ok := removeSet[p.ID]; !ok {
			newProdsList = append(newProdsList, p)
			order.Price = order.Price.Add(p.Price)
		}
	}

	order.Products = newProdsList

	return o.ordersRepo.UpdateOrder(ctx, order)
}

func (o *ordersUseCase) DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	// Check ownership
	_, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	return o.ordersRepo.DeleteOrder(ctx, orderID)
}

func (o *ordersUseCase) SubscribeToPaymentStatusUpdates() {
	for notification := range domain.PaymentStatusChannel {
		order, err := o.GetOrderByPaymentID(context.Background(), notification.PaymentID)
		if err != nil {
			return
		}

		status := domain.OrderStatusFromNotification(notification.Status)
		order.Status = status

		_, err = o.ordersRepo.UpdateOrder(context.Background(), order)
		if err != nil {
			o.logger.Errorw(
				"failed updating order",
				zap.Error(err),
			)
		}
		return
	}
}
