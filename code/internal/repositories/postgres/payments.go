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

func (p paymentsRepositoryImpl) PayOrder(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	return payment, nil
}

func NewPaymentsRepositoryImpl(log *zap.SugaredLogger, db *gorm.DB) ports.PaymentRepository {
	return &paymentsRepositoryImpl{log: log, db: db}
}

//const paymentsTable = "lanchonete_payments"
