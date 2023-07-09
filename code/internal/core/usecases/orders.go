package usecases

import (
	"context"

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
}

func (o *ordersUseCase) ListOrders(ctx context.Context, limit, offset int, userID uuid.UUID) (*domain.OrderList, error) {
	if !isAdmin(o.logger, o.userUC, ctx, userID) {
		return o.ordersRepo.ListOrdersByUser(ctx, limit, offset, userID)
	}
	return o.ordersRepo.ListOrders(ctx, limit, offset)
}

func NewOrdersUseCase(logger *zap.SugaredLogger, ordersRepo ports.OrdersRepository, userUC ports.UsersUseCase, prodUC ports.ProductsUseCase) ports.OrdersUseCase {
	return &ordersUseCase{logger: logger, ordersRepo: ordersRepo, userUC: userUC, prodUC: prodUC}
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

func (o *ordersUseCase) CreateOrder(ctx context.Context, userID uuid.UUID, products []uuid.UUID) (*domain.Order, error) {
	prodSum, err := o.prodUC.GetProductsPriceSumByID(ctx, products)
	if err != nil {
		return nil, err
	}

	var order *domain.Order
	order = domain.NewOrder(uuid.New(), userID, time.Now(), products)
	order.Price = prodSum.Sum

	return o.ordersRepo.CreateOrder(ctx, order)
}

func (o *ordersUseCase) InsertProductsIntoOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) (*domain.Order, error) {
	order, err := o.GetOrder(ctx, userID, orderID)
	// Check ownership
	if err != nil {
		return nil, err
	}

	// append products
	prodIDs := make([]uuid.UUID, 0, len(order.ProductsIDs)+len(products))
	for _, v := range order.ProductsIDs {
		prodIDs = append(prodIDs, v)
	}
	for _, v := range products {
		prodIDs = append(prodIDs, v)
	}

	err = o.getTotalAndUpdate(ctx, prodIDs, order)
	if err != nil {
		return nil, err
	}

	return order, err
}

func (o *ordersUseCase) RemoveProductFromOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) (*domain.Order, error) {
	order, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}

	// append products
	prodIDs := make([]uuid.UUID, 0, len(order.ProductsIDs))
	for _, v := range order.ProductsIDs {
		prodIDs = append(prodIDs, v)
	}
	for _, v := range products {
		removeProduct(prodIDs, v)
	}

	err = o.getTotalAndUpdate(ctx, prodIDs, order)
	if err != nil {
		return nil, err
	}

	return order, err
}

func (o *ordersUseCase) DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	// Check ownership
	_, err := o.GetOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	return o.ordersRepo.DeleteOrder(ctx, orderID)
}

func (o *ordersUseCase) SetOrderAsPaid(ctx context.Context, payment *domain.Payment) error {
	// Check ownership
	order, err := o.GetOrder(ctx, payment.UserID, payment.OrderID)
	if err != nil {
		return err
	}
	if order.Status == OrderPaidStatus {
		return nil
	}

	return o.ordersRepo.SetOrderAsPaid(ctx, payment)
}

func removeProduct(current []uuid.UUID, prodID uuid.UUID) []uuid.UUID {
	if len(current) == 0 {
		return make([]uuid.UUID, 0)
	}

	for k, v := range current {
		if v == prodID {
			return append(current[:k], current[k+1:]...)
		}
	}

	return current
}

func (o *ordersUseCase) getTotalAndUpdate(ctx context.Context, prodIDs []uuid.UUID, order *domain.Order) error {
	// get value
	sum, err := o.prodUC.GetProductsPriceSumByID(ctx, prodIDs)
	if err != nil {
		return err
	}

	order.ProductsIDs = prodIDs
	order.Price = sum.Sum

	//update order
	updated, err := o.ordersRepo.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	order = updated
	return nil
}
