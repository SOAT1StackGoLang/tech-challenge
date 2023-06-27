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
	log *zap.SugaredLogger
	db  *gorm.DB
}

func (o *ordersRepositoryImpl) CreateOrder(ctx context.Context, userID, productID uuid.UUID) (*domain.Order, error) {
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
