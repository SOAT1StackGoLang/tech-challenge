package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ordersRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func (o *ordersRepositoryImpl) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*domain.Order, error) {
	order := &Order{}

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("*").
		Where("payment_id = ?", paymentID).
		First(order).Error; err != nil {
		o.log.Errorw(
			"db failed getting order",
			zap.Error(err),
		)
		return nil, err
	}

	out := order.toDomain()
	return out, err
}

func (o *ordersRepositoryImpl) ListOrdersByUser(ctx context.Context, limit, offset int, userID uuid.UUID) (*domain.OrderList, error) {
	var orders []Order
	var total int64

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&orders).Error; err != nil {
		o.log.Errorw(
			"failed listing orders",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		o.log.Errorw(
			"failed counting orders by user_id",
			zap.String("category", userID.String()),
			zap.Error(err),
		)
	}

	oList := &domain.OrderList{}
	out := make([]*domain.Order, 0, len(orders))

	for _, v := range orders {
		out = append(out, v.toDomain())
	}

	oList.Orders = out
	oList.Total = total
	oList.Limit = limit
	oList.Offset = offset

	return oList, err

}

func (o *ordersRepositoryImpl) ListOrders(ctx context.Context, limit, offset int) (*domain.OrderList, error) {
	var total int64

	var saveOrders []Order

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Limit(limit).
		Offset(offset).
		Order("status DESC").
		Where("status > ? AND status < ? ", ORDER_STATUS_WAITING_PAYMENT, ORDER_STATUS_FINISHED).
		Scan(&saveOrders).Error; err != nil {
		o.log.Errorw(
			"failed listing orders",
			zap.Error(err),
		)
		return nil, err
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Where("status > ? AND status < ? ", ORDER_STATUS_WAITING_PAYMENT, ORDER_STATUS_FINISHED).
		Count(&total).Error; err != nil {
		o.log.Errorw(
			"failed counting orders",
			zap.Error(err),
		)
	}

	oList := &domain.OrderList{}
	out := make([]*domain.Order, 0, len(saveOrders))

	for _, v := range saveOrders {
		out = append(out, v.toDomain())
	}

	oList.Orders = out
	oList.Total = total
	oList.Limit = limit
	oList.Offset = offset

	return oList, err
}

const ordersTable = "lanchonete_orders"

func NewPgxOrdersRepository(log *zap.SugaredLogger, db *gorm.DB) ports.OrdersRepository {
	return &ordersRepositoryImpl{log: log, db: db}
}

func (o *ordersRepositoryImpl) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	order := &Order{}

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("*").
		Where("id = ?", orderID).
		First(order).Error; err != nil {
		o.log.Errorw(
			"db failed getting order",
			zap.Error(err),
		)
		return nil, err
	}

	out := order.toDomain()
	return out, err
}

func (o *ordersRepositoryImpl) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	var out *domain.Order

	in := Order{}
	in.fromDomain(order)
	in.Status = ORDER_STATUS_OPEN

	if err := o.db.WithContext(ctx).Table(ordersTable).Omit("updated_at").Create(&in).Error; err != nil {
		o.log.Errorw(
			"db failed at CreateOrder",
			zap.Any("order_input", order),
			zap.Error(err),
		)
		return nil, err
	}

	out = in.toDomain()

	return out, nil
}

func (o *ordersRepositoryImpl) UpdateOrder(ctx context.Context, in *domain.Order) (*domain.Order, error) {
	order := &Order{}
	order.fromDomain(in)

	order.UpdatedAt = sql.NullTime{
		Time:  in.UpdatedAt,
		Valid: true,
	}

	var oS OrderStatus
	order.Status = oS.fromDomain(in.Status)

	if err := o.db.WithContext(ctx).Table(ordersTable).
		Updates(&order).
		Where("id = ?", in.ID).
		Error; err != nil {
		o.log.Errorw(
			"db failed updating order",
			zap.Any("in_order", in),
			zap.Any("repo_order", order),
			zap.Error(err),
		)
		return nil, err
	}

	return order.toDomain(), nil
}

func (o *ordersRepositoryImpl) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	deletedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err := o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumn("deleted_at", deletedAt).
		Where("order_id", orderID).
		Error; err != nil {
		o.log.Errorw(
			"db failed deleting order",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
	}
	return nil
}

func (o *ordersRepositoryImpl) SetOrderAsPaid(ctx context.Context, payment *domain.Payment) (err error) {
	//updatedAt := sql.NullTime{
	//	Time:  time.Now(),
	//	Valid: true,
	//}

	//if err = o.db.WithContext(ctx).Table(ordersTable).
	//	Where("id = ?", payment.OrderID).
	//	UpdateColumns(map[string]any{
	//		"status":     usecases.OrderPaidStatus,
	//		"updated_at": updatedAt,
	//		"payment_id": payment.ID,
	//	}).
	//	Error; err != nil {
	//	o.log.Errorw(
	//		"db failed finishing order",
	//		zap.String("order_id", payment.OrderID.String()),
	//		zap.Error(err),
	//	)
	//}

	return err
}
