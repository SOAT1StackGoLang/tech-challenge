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
	logger      *zap.SugaredLogger
	orderUC     ports.OrdersUseCase
	paymentRepo ports.PaymentRepository
}

const OrderPaidStatus = "PAGA"

func NewPaymentsUseCase(logger *zap.SugaredLogger, orderUC ports.OrdersUseCase, repo ports.PaymentRepository) ports.PaymentUseCase {
	return &paymentsUseCase{logger: logger, orderUC: orderUC, paymentRepo: repo}
}

func (p *paymentsUseCase) PayOrder(ctx context.Context, orderID, userID uuid.UUID) (*domain.Payment, error) {
	payment := domain.NewPayment(uuid.New(), time.Now(), orderID, userID)

	receipt, err := p.paymentRepo.PayOrder(ctx, payment)
	if err != nil {
		return nil, err
	}

	if err = p.orderUC.SetOrderAsPaid(ctx, receipt); err != nil {
		return nil, err
	}

	return receipt, nil
}
