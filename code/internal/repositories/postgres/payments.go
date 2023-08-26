package postgres

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type paymentsRepositoryImpl struct {
	log     *zap.SugaredLogger
	db      *gorm.DB
	orderUC ports.OrdersUseCase
}

const paymentTable = "lanchonete_payments"

func (p paymentsRepositoryImpl) CreatePayment(ctx context.Context, in *domain.Payment) (*domain.Payment, error) {
	payment := new(Payment)
	payment.fromDomain(in)

	if err := p.db.WithContext(ctx).Table(paymentTable).Create(&payment).Error; err != nil {
		p.log.Errorw(
			"db failed at CreatePayment",
			zap.Any("payment_input", in),
			zap.Error(err),
		)
		return nil, err
	}

	out := new(domain.Payment)
	out = payment.toDomain()

	return out, nil
}

func NewPaymentsRepository(log *zap.SugaredLogger, db *gorm.DB) ports.PaymentRepository {
	return &paymentsRepositoryImpl{log: log, db: db}
}
