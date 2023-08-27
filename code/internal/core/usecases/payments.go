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
	paymentRepo ports.PaymentRepository
}

func (p *paymentsUseCase) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	return p.paymentRepo.GetPayment(ctx, paymentID)
}

func NewPaymentsUseCase(logger *zap.SugaredLogger, repo ports.PaymentRepository) ports.PaymentUseCase {
	return &paymentsUseCase{logger: logger, paymentRepo: repo}
}

func (p *paymentsUseCase) CreatePayment(ctx context.Context, order *domain.Order) (*domain.Payment, error) {

	payment := domain.NewPayment(uuid.New(), time.Now(), order.ID, order.Price, domain.PAYMENT_STATUS_OPEN)

	receipt, err := p.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (p *paymentsUseCase) UpdatePayment(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) (*domain.Payment, error) {
	payment, err := p.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	payment.UpdatedAt = time.Now()
	payment.Status = status

	updated, err := p.paymentRepo.UpdatePayment(ctx, payment)

	defer func() {
		p.PublishPaymentStatus(domain.PaymentStatusNotification{
			PaymentID: payment.ID,
			Status:    status,
			OrderID:   payment.OrderID,
		})
	}()

	return updated, err
}

func (p *paymentsUseCase) PublishPaymentStatus(notification domain.PaymentStatusNotification) {
	domain.PaymentStatusChannel <- notification
}
