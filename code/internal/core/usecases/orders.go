package usecases

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ordersUseCase struct {
	logger     *zap.SugaredLogger
	ordersRepo ports.OrdersRepository
	userUC     ports.UsersUseCase
}

func NewOrdersUseCase(logger *zap.SugaredLogger, ordersRepo ports.OrdersRepository, userUC ports.UsersUseCase) ports.OrdersUseCase {
	return &ordersUseCase{logger: logger, ordersRepo: ordersRepo, userUC: userUC}
}

func (o ordersUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return o.ordersRepo.GetOrder(ctx, orderID)
}

func (o ordersUseCase) CreateOrder(ctx context.Context, userID uuid.UUID, products []uuid.UUID) (*domain.Order, error) {
	return o.ordersRepo.CreateOrder(ctx, userID, products)
}

func (o ordersUseCase) InsertProductsIntoOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) error {
	return o.InsertProductsIntoOrder(ctx, userID, orderID, products)
}

func (o ordersUseCase) RemoveProductFromOrder(ctx context.Context, userID, orderID uuid.UUID, products []uuid.UUID) error {
	return o.ordersRepo.RemoveProductsFromOrder(ctx, userID, orderID, products)
}

func (o ordersUseCase) DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	return o.ordersRepo.DeleteOrder(ctx, userID, orderID)
}

func (o ordersUseCase) PayOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	err := validateIsAdmin(o.logger, o.userUC, ctx, userID)
	if err != nil {
		return err
	}
	return o.ordersRepo.FinishOrder(ctx, orderID)
}
