package postgres

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	out.Products = make([]domain.Product, 0, len(in.Products))
	for _, v := range in.Products {
		product, err := o.prodUsc.GetProduct(ctx, v)
		if err != nil {
			o.log.Errorw(
				"order db failed getting product by id",
				zap.String("product_id", v.String()),
				zap.Error(err),
			)
			return nil, err
		}

		out.Products = append(out.Products, *product)
	}
	out = in.toDomain(out.Products)

	return out, err

	//TODO implement me
	panic("implement me")
}

func (o *ordersRepositoryImpl) InsertProductIntoOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersRepositoryImpl) RemoveProductFromOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordersRepositoryImpl) DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewPgxOrdersRepository(log *zap.SugaredLogger, db *gorm.DB) ports.OrdersRepository {
	return &ordersRepositoryImpl{log: log, db: db}
}
