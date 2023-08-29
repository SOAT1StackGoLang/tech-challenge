package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

type PaymentsHttpHandler struct {
	ctx             context.Context
	paymentsUseCase ports.PaymentUseCase
	ordersUseCase   ports.OrdersUseCase
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
