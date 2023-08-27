package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"time"
)

type (
	PaymentsHttpHandler struct {
		ctx             context.Context
		paymentsUseCase ports.PaymentUseCase
		ordersUseCase   ports.OrdersUseCase
	}

	PaymentNotification struct {
		ID       string `json:"id" description:"ID do pagamento"`
		OrderID  string `json:"order_id" description:"ID do pedido a ser pago"`
		Approved bool   `json:"approved" description:"True para aprovado false para recusado"`
	}

	PaymentStatus string

	Payment struct {
		ID        string        `json:"id" description:"ID do pagamento"`
		OrderID   string        `json:"order_id" description:"ID do pedido a ser pago"`
		CreatedAt string        `json:"created_at" description:"Data de criação"`
		UpdatedAt string        `json:"updated_at" description:"Data de Atualização"`
		Value     string        `json:"value" description:"Valor em R$"`
		Status    PaymentStatus `json:"status do pagamento"`
	}
)

func (pN *PaymentNotification) toDomain() *domain.PaymentStatusNotification {
	var pS PaymentStatus
	pS = pS.fromRequest(pN.Approved)
	return &domain.PaymentStatusNotification{
		PaymentID: helpers.SafeUUIDFromString(pN.ID),
		OrderID:   helpers.SafeUUIDFromString(pN.OrderID),
		Status:    pS.toDomain(),
	}
}

func (p *Payment) fromDomain(payment *domain.Payment) {
	p.ID = payment.ID.String()
	p.OrderID = payment.OrderID.String()
	p.CreatedAt = payment.CreatedAt.Format(time.RFC3339)

	if !payment.UpdatedAt.IsZero() {
		p.UpdatedAt = payment.UpdatedAt.Format(time.RFC3339)
	}
	var price string
	if !payment.Price.IsZero() {
		price = helpers.ParseDecimalToString(payment.Price)
		p.Value = price
	}

	var pS PaymentStatus
	p.Status = pS.fromDomain(payment.Status)
}

const (
	PAYMENT_STATUS_OPEN     PaymentStatus = "Aguardando Pagamento"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_STATUS_REFUSED                = "Recusado"
)

func (pS PaymentStatus) toDomain() domain.PaymentStatus {
	switch pS {
	case PAYMENT_STATUS_OPEN:
		return domain.PAYMENT_STATUS_OPEN
	case PAYMENT_STATUS_APPROVED:
		return domain.PAYMENT_STATUS_APPROVED
	}
	return domain.PAYMENT_SATUS_REFUSED
}

func (pS PaymentStatus) fromDomain(status domain.PaymentStatus) PaymentStatus {
	switch status {
	case domain.PAYMENT_STATUS_OPEN:
		return PAYMENT_STATUS_OPEN
	case domain.PAYMENT_STATUS_APPROVED:
		return PAYMENT_STATUS_APPROVED
	}
	return PAYMENT_STATUS_REFUSED
}

func (pS PaymentStatus) fromRequest(approved bool) PaymentStatus {
	if approved {
		return PAYMENT_STATUS_APPROVED
	}
	return PAYMENT_STATUS_REFUSED
}

func NewPaymentsHttpHandler(
	ctx context.Context,
	paymentsUseCase ports.PaymentUseCase,
	ws *restful.WebService,
) *PaymentsHttpHandler {
	handler := &PaymentsHttpHandler{ctx: ctx, paymentsUseCase: paymentsUseCase}

	tags := []string{"payments"}

	ws.Route(ws.POST("/webhook/payment-notification").To(handler.handlePaymentNotification).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Efetua pagamento de pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(PaymentNotification{}).
		Returns(http.StatusOK, "Pagamento efetuado com sucesso", Payment{}).
		Returns(http.StatusBadRequest, "Requisição incorreta", nil).
		Returns(http.StatusInternalServerError, "Falha do servidor", nil))

	return handler
}

func (pHH *PaymentsHttpHandler) handlePaymentNotification(request *restful.Request, response *restful.Response) {
	var pN PaymentNotification
	if err := request.ReadEntity(&pN); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	notification := pN.toDomain()

	p, err := pHH.paymentsUseCase.UpdatePayment(pHH.ctx, notification.PaymentID, notification.Status)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	out := new(Payment)
	out.fromDomain(p)

	_ = response.WriteAsJson(out)
}
