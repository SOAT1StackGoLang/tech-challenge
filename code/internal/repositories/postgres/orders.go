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
	log     *zap.SugaredLogger
	db      *gorm.DB
	prodUsc ports.ProductsUseCase
}

const ordersTable = "lanchonete_orders"

func (o *ordersRepositoryImpl) CreateOrder(ctx context.Context, userID, productID uuid.UUID) (*domain.Order, error) {
	var out *domain.Order

	var in *Order
	in.newFromDomain(userID, productID)

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).Create(&in).Error; err != nil {
		o.log.Errorw(
			"db failed inserting order",
			zap.Error(err),
		)
		return nil, err
	}

	if err = o.productsIDToDomainProducts(ctx, out, in); err != nil {
		return nil, err
	}

	out = in.toDomain(out.Products)

	return out, err
}

func (o *ordersRepositoryImpl) productsIDToDomainProducts(ctx context.Context, out *domain.Order, in *Order) error {
	out.Products = make([]domain.Product, 0, len(in.Products))
	for _, v := range in.Products {
		product, err := o.prodUsc.GetProduct(ctx, v)
		if err != nil {
			o.log.Errorw(
				"order db failed getting product by id",
				zap.String("product_id", v.String()),
				zap.Error(err),
			)
			return err
		}

		out.Products = append(out.Products, *product)
	}
	return nil
}

func (o *ordersRepositoryImpl) InsertProductIntoOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error) {
	var out *domain.Order

	var in *Order

	var err error
	if err = o.db.WithContext(ctx).Table(ordersTable).
		Select("*").
		Where("id = ?", orderID).
		Where("user_id = ?", userID).
		Where("deleted_at = ?", nil).
		First(in).
		Error; err != nil {
		o.log.Errorw(
			"db failed getting order by id and user id",
			zap.String("order_id", orderID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	in.Products = append(in.Products, productID)
	in.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = o.db.WithContext(ctx).Table(ordersTable).
		Updates(in).
		Where("id = ?", in.ID).
		Error; err != nil {
		o.log.Errorw(
			"db failed inserting product into order",
			zap.String("order_id", orderID.String()),
			zap.String("product_id", productID.String()),
		)
	}

	if err = o.productsIDToDomainProducts(ctx, out, in); err != nil {
		return nil, err
	}
	out = in.toDomain(out.Products)

	return out, err
}

func (o *ordersRepositoryImpl) RemoveProductFromOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
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
			zap.String("order_id", orderID.String()))
		return err
	}
	return err
}

func NewPgxOrdersRepository(log *zap.SugaredLogger, db *gorm.DB, productUseCase ports.ProductsUseCase) ports.OrdersRepository {
	return &ordersRepositoryImpl{log: log, db: db, prodUsc: productUseCase}
}
