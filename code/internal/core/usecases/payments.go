package usecases

import (
	"context"
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

func (p *paymentsUseCase) PayOrder(ctx context.Context, orderID, userID uuid.UUID) (*domain.Payment, error) {
	payment := domain.NewPayment(uuid.New(), time.Now(), orderID, userID)

	if err := p.orderUC.SetOrderAsPaid(ctx, userID, orderID); err != nil {
		return nil, err
	}

	return payment, nil
}
