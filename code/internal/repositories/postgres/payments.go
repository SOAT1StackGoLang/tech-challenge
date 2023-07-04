package postgres

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type paymentsRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func (p paymentsRepositoryImpl) PayOrder(ctx context.Context, userID, orderID uuid.UUID) error {
	var payment Payment
	payment.newPayment(userID, orderID)
	//TODO implement me
	panic("implement me")
}

func NewPaymentsRepositoryImpl(log *zap.SugaredLogger, db *gorm.DB) ports.PaymentRepository {
	return &paymentsRepositoryImpl{log: log, db: db}
}

const paymentsTable = "lanchonete_payments"
