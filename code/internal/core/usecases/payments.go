package usecases

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type paymentsUseCase struct {
	logger  *zap.SugaredLogger
	orderUC ports.OrdersUseCase
}

func NewPaymentsUseCase(logger *zap.SugaredLogger, orderUC ports.OrdersUseCase) ports.PaymentUseCase {
	return &paymentsUseCase{logger: logger, orderUC: orderUC}
}

func (p paymentsUseCase) PayOrder(ctx context.Context, orderID, userID uuid.UUID) (*domain.Payment, error) {
	payment := domain.NewPayment(uuid.New(), time.Now(), orderID, userID)

	order, err := p.orderUC.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, helpers.ErrUnauthorized
	}

	if err = p.orderUC.PayOrder(ctx, userID, orderID); err != nil {
		return nil, err
	}

	return payment, err
}
