package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
)

type (
	PaymentsHttpHandler struct {
		ctx             context.Context
		paymentsUseCase ports.PaymentUseCase
	}

	PaymentRequest struct {
		OrderID string `json:"order_id" description:"ID do pedido a ser pago"`
		UserID  string `json:"user_id" description:"Usuário efetuando o pagamento"`
	}

	Payment struct {
		PaymentID string `json:"payment_id" description:"ID do pagamento"`
		PaymentRequest
	}
)

func (p *Payment) fromDomain(payment *domain.Payment) {
	if p == nil {
		p = &Payment{}
	}

	p.PaymentID = payment.ID.String()

	var pR PaymentRequest
	pR.fromDomain(payment)
	p.PaymentRequest = pR
}

func (pR *PaymentRequest) fromDomain(payment *domain.Payment) {
	if pR == nil {
		pR = &PaymentRequest{}
	}

	pR.OrderID = payment.OrderID.String()
	pR.UserID = payment.UserID.String()
}

func NewPaymentsHttpHandler(ctx context.Context, paymentsUseCase ports.PaymentUseCase, ws *restful.WebService) *PaymentsHttpHandler {
	handler := &PaymentsHttpHandler{ctx: ctx, paymentsUseCase: paymentsUseCase}

	tags := []string{"payments"}

	ws.Route(ws.POST("/payments").To(handler.PayOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Efetua pagamento de pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(PaymentRequest{}).
		Returns(http.StatusOK, "Pagamento efetuado com sucesso", Payment{}).
		Returns(http.StatusBadRequest, "Requisição incorreta", nil).
		Returns(http.StatusInternalServerError, "Falha do servidor", nil))

	return handler
}

func (pH *PaymentsHttpHandler) PayOrder(request *restful.Request, response *restful.Response) {
	var pR PaymentRequest

	if err := request.ReadEntity(&pR); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(pR.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	oID, err := uuid.Parse(pR.OrderID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	payment, err := pH.paymentsUseCase.PayOrder(pH.ctx, oID, uid)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var p Payment
	p.fromDomain(payment)
	_ = response.WriteAsJson(p)

}
