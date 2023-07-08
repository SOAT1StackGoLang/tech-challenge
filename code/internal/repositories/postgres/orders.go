package postgres

import (
	"context"
	"database/sql"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ordersRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

const ordersTable = "lanchonete_orders"

func NewPgxOrdersRepository(log *zap.SugaredLogger, db *gorm.DB) ports.OrdersRepository {
	return &ordersRepositoryImpl{log: log, db: db}
}

func (o *ordersRepositoryImpl) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	var order *Order

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

	var in *Order
	in.newFromDomain(order)

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).Create(&in).Error; err != nil {
		o.log.Errorw(
			"db failed creating order",
			zap.Error(err),
		)
		return nil, err
	}

	out = in.toDomain()

	return out, err
}

func (o *ordersRepositoryImpl) UpdateOrder(ctx context.Context, in *domain.Order) (*domain.Order, error) {
	order := &Order{}
	order.fromDomain(in)

	order.UpdatedAt.Time = time.Now()

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
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

	return order.toDomain(), err
}

func (o *ordersRepositoryImpl) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	var err error
	deletedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err = o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumn("deleted_at", deletedAt).
		Where("order_id", orderID).
		Error; err != nil {
		o.log.Errorw(
			"db failed deleting order",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
	}
	return err
}

func (o *ordersRepositoryImpl) SetOrderAsPaid(ctx context.Context, payment *domain.Payment) (err error) {
	status := "PAGA"
	updatedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumns(map[string]any{
			"status":     status,
			"updated_at": updatedAt,
			"payment_id": payment.ID,
		}).
		Where("id = ?", payment.OrderID).
		Error; err != nil {
		o.log.Errorw(
			"db failed finishing order",
			zap.String("order_id", payment.OrderID.String()),
			zap.Error(err),
		)
	}

	return err
}
