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

func NewPaymentsUseCase(logger *zap.SugaredLogger, repo ports.PaymentRepository) ports.PaymentUseCase {
	return &paymentsUseCase{logger: logger, paymentRepo: repo}
}

func (p *paymentsUseCase) CreatePayment(ctx context.Context, order *domain.Order) (*domain.Payment, error) {

	payment := domain.NewPayment(uuid.New(), time.Now(), order.ID, order.Price)

	receipt, err := p.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
