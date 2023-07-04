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

func (o *ordersRepositoryImpl) CreateOrder(ctx context.Context, userID uuid.UUID, products []uuid.UUID) (*domain.Order, error) {
	var out *domain.Order

	var in *Order
	in.newFromDomain(userID, products)

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).Create(&in).Error; err != nil {
		o.log.Errorw(
			"db failed creating order",
			zap.Error(err),
		)
		return nil, err
	}

	out = in.toDomain(products)

	return out, err
}

func (o *ordersRepositoryImpl) InsertProductsIntoOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) (err error) {
	var currentProducts []uuid.UUID

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("products").
		Where("id = ?", orderID).
		Where("user_id = ?", userID).
		Where("deleted_at = ?", nil).
		First(&currentProducts).
		Error; err != nil {
		o.log.Errorw(
			"db failed getting order by id and user id",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
		return err
	}

	updatedProducts := currentProducts
	for _, p := range products {
		updatedProducts = append(updatedProducts, p)
	}

	updatedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Updates(map[string]any{
			"updated_at": updatedAt,
			"products":   updatedProducts,
		}).
		Where("id = ?", orderID).
		Error; err != nil {
		o.log.Errorw(
			"db failed inserting product into order",
			zap.String("order_id", orderID.String()),
			zap.Any("products", products),
		)
	}

	return err
}

func (o *ordersRepositoryImpl) RemoveProductsFromOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) (err error) {
	var currentProducts []uuid.UUID

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("products").
		Where("id = ?", orderID).
		Where("user_id = ?", userID).
		Where("deleted_at = ?", nil).
		First(&currentProducts).
		Error; err != nil {
		o.log.Errorw(
			"db failed getting order by id and user id",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
		return err
	}

	var updatedProducts []uuid.UUID
	for _, p := range products {
		updatedProducts = removeProduct(currentProducts, p)
	}

	if len(updatedProducts) == 0 {
		o.log.Debugw(
			"order has no more items, deleting",
			zap.String("order_id", orderID.String()),
		)
		return o.DeleteOrder(ctx, userID, orderID)
	}

	updatedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Updates(map[string]any{
			"products":   updatedProducts,
			"updated_at": updatedAt,
		}).
		Where("id = ?", orderID).
		Error; err != nil {
		o.log.Errorw(
			"db failed updating order with removed item",
			zap.String("order_id", orderID.String()),
			zap.Any("products", updatedProducts),
			zap.Error(err),
		)
	}

	return err
}

func (o *ordersRepositoryImpl) DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	var err error
	deletedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err = o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumn("deleted_at", deletedAt).
		Where("user_id", userID).
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

func (o *ordersRepositoryImpl) FinishOrder(ctx context.Context, orderID uuid.UUID) (err error) {
	status := "DONE"
	updatedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		UpdateColumns(map[string]any{
			"status":     status,
			"updated_at": updatedAt,
		}).
		Where("id = ?").
		Error; err != nil {
		o.log.Errorw(
			"db failed finishing order",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
	}

	return err
}

func removeProduct(current []uuid.UUID, prodID uuid.UUID) []uuid.UUID {
	if len(current) == 0 {
		return nil
	}

	for k, v := range current {
		if v == prodID {
			return append(current[:k], current[k+1:]...)
		}
	}

	return current
}
