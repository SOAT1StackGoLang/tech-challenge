package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
)

type (
	OrdersHttpHandler struct {
		ctx      context.Context
		ordersUC ports.OrdersUseCase
	}

	Checkout struct {
		Order       Order       `json:"order" description:"Pedido"`
		PaymentInfo PaymentInfo `json:"payment_info" description:"Informações de Cobrança"`
	}

	PaymentInfo struct {
		PaymentID string `json:"payment_id"`
		Value     string `json:"value" description:"Valor a ser pago"`
	}

	OrderStatus string

	Order struct {
		ID        uuid.UUID   `json:"id" description:"ID do Pedido"`
		PaymentID string      `json:"payment_id,omitempty" description:"ID do pagamento"`
		CreatedAt string      `json:"created_at" description:"Data de criação"`
		UpdatedAt string      `json:"updated_at,omitempty" description:"Data de atualização"`
		DeletedAt string      `json:"deleted_at,omitempty" description:"Data de deleção"`
		Price     string      `json:"price" description:"Preço do pedido"`
		Status    OrderStatus `json:"status" description:"Status do pedido"`
		Products  []Product   `json:"products" description:"Lista de Pedidos"`
	}

	InsertionOrder struct {
		UserID      string      `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []uuid.UUID `json:"products_ids" description:"ID dos produtos"`
	}

	InsertionOrderSwagger struct {
		UserID      string   `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []string `json:"products_ids" description:"Lista de ID dos produtos separados por vírgula"`
	}

	UpdateOrder struct {
		ID string `json:"id" description:"ID do Pedido"`
		InsertionOrder
	}

	OrderList struct {
		Orders []Order `json:"orders"`
		Limit  int     `json:"limit" default:"10"`
		Offset int     `json:"offset"`
		Total  int64   `json:"total"`
	}

	OrderListRequest struct {
		UserID string `json:"user_id"`
		Limit  int    `json:"limit" default:"10" description:"Quantidade de registros"`
		Offset int    `json:"offset"`
	}

	OrderCheckoutRequest struct {
		UserID  string `json:"user_id"`
		OrderID string `json:"order_id" description:"ID do Pedido"`
	}

	OrderStatusUpdate struct {
		OrderID string `json:"order_id" description:"Código de identificação do pedido"`
		UserID  string `json:"user_id" description:"Código de descrição do usuário requerente"`
		Status  string `json:"status" description:"Status para qual deseja mudar o pedido" enum:"Recebido|Preparacao|Pronto|Finalizado|Cancelado"`
	}
)

func (o *Order) fromDomain(order *domain.Order) {
	if o == nil {
		o = &Order{}
	}

	o.ID = order.ID

	var p string
	if !order.Price.IsZero() {
		p = helpers.ParseDecimalToString(order.Price)
		o.Price = p
	}

	var products []Product
	for _, dP := range order.Products {
		p := Product{}
		p.fromDomain(&dP)
		products = append(products, p)
	}

	o.Products = products

	if order.PaymentID != uuid.Nil {
		o.PaymentID = order.PaymentID.String()
	}
	o.CreatedAt = order.CreatedAt.Format(time.RFC3339)
	if !order.UpdatedAt.IsZero() {
		o.UpdatedAt = order.UpdatedAt.Format(time.RFC3339)
	} else {
		o.UpdatedAt = ""
	}
	if !order.DeletedAt.IsZero() {
		o.DeletedAt = order.DeletedAt.Format(time.RFC3339)
	} else {
		o.DeletedAt = ""
	}

	orderStatus := new(OrderStatus)
	o.Status = orderStatus.fromDomain(order.Status)
}

const (
	ORDER_STATUS_UNSET           OrderStatus = ""
	ORDER_STATUS_OPEN                        = "Aberto"
	ORDER_STATUS_WAITING_PAYMENT             = "Aguardando Pagamento"
	ORDER_STATUS_RECEIVED                    = "Recebido"
	ORDER_STATUS_PREPARING                   = "Preparacao"
	ORDER_STATUS_DONE                        = "Pronto"
	ORDER_STATUS_FINISHED                    = "Finalizado"
	ORDER_STATUS_CANCELED                    = "Cancelado"
)

func (oS *OrderStatus) fromDomain(status domain.OrderStatus) OrderStatus {
	switch status {
	case domain.ORDER_STATUS_UNSET:
		return ORDER_STATUS_UNSET
	case domain.ORDER_STATUS_OPEN:
		return ORDER_STATUS_OPEN
	case domain.ORDER_STATUS_WAITING_PAYMENT:
		return ORDER_STATUS_WAITING_PAYMENT
	case domain.ORDER_STATUS_RECEIVED:
		return ORDER_STATUS_RECEIVED
	case domain.ORDER_STATUS_PREPARING:
		return ORDER_STATUS_PREPARING
	case domain.ORDER_STATUS_DONE:
		return ORDER_STATUS_DONE
	case domain.ORDER_STATUS_FINISHED:
		return ORDER_STATUS_FINISHED
	}
	return ORDER_STATUS_CANCELED
}

func stringToDomainStatus(status string) domain.OrderStatus {
	switch status {
	case ORDER_STATUS_RECEIVED:
		return domain.ORDER_STATUS_RECEIVED
	case ORDER_STATUS_PREPARING:
		return domain.ORDER_STATUS_PREPARING
	case ORDER_STATUS_DONE:
		return domain.ORDER_STATUS_DONE
	case ORDER_STATUS_FINISHED:
		return domain.ORDER_STATUS_FINISHED
	}
	return domain.ORDER_STATUS_UNSET
}

func (oH *OrdersHttpHandler) handleGetOrder(request *restful.Request, response *restful.Response) {
	var queryStruct QueryStruct

	if err := request.ReadEntity(&queryStruct); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id, uid, err := queryStruct.parseToUuid()
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
	}

	order, err := oH.ordersUC.GetOrder(oH.ctx, uid, id)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleCreateOrder(request *restful.Request, response *restful.Response) {
	var insertOrder InsertionOrder
	if err := request.ReadEntity(&insertOrder); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	dPs := oH.productsToDomainProducts(insertOrder.ProductsIDs)

	order, err := oH.ordersUC.CreateOrder(oH.ctx, helpers.SafeUUIDFromString(insertOrder.UserID), dPs)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)

		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) productsToDomainProducts(products []uuid.UUID) []domain.Product {
	var dPs []domain.Product
	for _, p := range products {
		dP := &domain.Product{ID: p}

		dPs = append(dPs, *dP)
	}
	return dPs
}

func (oH *OrdersHttpHandler) handleAddProductsIntoOrder(request *restful.Request, response *restful.Response) {
	var addRequest UpdateOrder
	if err := request.ReadEntity(&addRequest); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id := helpers.SafeUUIDFromString(addRequest.ID)
	uid := helpers.SafeUUIDFromString(addRequest.InsertionOrder.UserID)

	dPs := oH.productsToDomainProducts(addRequest.InsertionOrder.ProductsIDs)

	order, err := oH.ordersUC.InsertProductsIntoOrder(oH.ctx, uid, id, dPs)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleStatusUpdate(request *restful.Request, response *restful.Response) {
	var req OrderStatusUpdate
	if err := request.ReadEntity(&req); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	oID := helpers.SafeUUIDFromString(req.OrderID)
	uID := helpers.SafeUUIDFromString(req.UserID)

	dS := stringToDomainStatus(req.Status)
	if dS == domain.ORDER_STATUS_UNSET {
		_ = response.WriteError(http.StatusBadRequest, errors.New("bad request: invalid status input"))
		return
	}

	resp, err := oH.ordersUC.UpdateOrderStatus(oH.ctx, uID, oID, dS)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Order
	out.fromDomain(resp)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleRemoveProductsOfOrder(request *restful.Request, response *restful.Response) {
	var removeReq UpdateOrder
	if err := request.ReadEntity(&removeReq); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id := helpers.SafeUUIDFromString(removeReq.ID)
	uid := helpers.SafeUUIDFromString(removeReq.InsertionOrder.UserID)

	dPs := oH.productsToDomainProducts(removeReq.ProductsIDs)
	order, err := oH.ordersUC.RemoveProductFromOrder(oH.ctx, uid, id, dPs)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleDeleteOrder(request *restful.Request, response *restful.Response) {
	var removeStruct QueryStruct

	if err := request.ReadEntity(&removeStruct); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id := helpers.SafeUUIDFromString(removeStruct.ID)
	uid := helpers.SafeUUIDFromString(removeStruct.UserID)

	err := oH.ordersUC.DeleteOrder(oH.ctx, uid, id)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (oH *OrdersHttpHandler) handleCheckout(request *restful.Request, response *restful.Response) {
	var oC OrderCheckoutRequest
	if err := request.ReadEntity(&oC); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(oC.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	id := helpers.SafeUUIDFromString(oC.OrderID)
	order, err := oH.ordersUC.Checkout(oH.ctx, uid, id)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Checkout
	var outOrder Order
	outOrder.fromDomain(order)

	var outPayment PaymentInfo
	outPayment.Value = outOrder.Price
	outPayment.PaymentID = outOrder.PaymentID

	out.Order = outOrder
	out.PaymentInfo = outPayment

	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleListOrders(request *restful.Request, response *restful.Response) {
	var oLR OrderListRequest
	if err := request.ReadEntity(&oLR); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(oLR.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	list, err := oH.ordersUC.ListOrders(oH.ctx, oLR.Limit, oLR.Offset, uid)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var oL OrderList
	var ord Order
	for _, v := range list.Orders {
		ord.fromDomain(v)
		oL.Orders = append(oL.Orders, ord)
	}
	oL.Total = list.Total
	oL.Limit = list.Limit
	oL.Offset = list.Offset

	_ = response.WriteAsJson(oL)
}

func NewOrdersHttpHandler(ctx context.Context, ordersUC ports.OrdersUseCase, ws *restful.WebService) *OrdersHttpHandler {
	handler := &OrdersHttpHandler{
		ctx:      ctx,
		ordersUC: ordersUC,
	}

	tags := []string{"orders"}

	ws.Route(ws.POST("/orders/get").To(handler.handleGetOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Obtém dados do pedido em função dos dados fornecidos").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(QueryStruct{}).
		Returns(http.StatusOK, "ok", Order{}).
		Returns(http.StatusBadRequest, "bad request", nil))
	ws.Route(ws.POST("/orders/all").To(handler.handleListOrders).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Lista pedidos").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(OrderListRequest{}).
		Returns(http.StatusOK, "sucesso", OrderList{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	ws.Route(ws.POST("/orders").To(handler.handleCreateOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionOrderSwagger{}).
		Returns(http.StatusOK, "sucesso", Order{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	ws.Route(ws.PUT("/orders/add").To(handler.handleAddProductsIntoOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Adiciona items ao pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateOrder{}).
		Returns(http.StatusOK, "sucesso", Order{}))
	ws.Route(ws.PUT("/orders/remove").To(handler.handleRemoveProductsOfOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Remove items do pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateOrder{}).
		Returns(http.StatusOK, "sucesso", Order{}).
		Returns(http.StatusBadRequest, "request incorreto", nil))
	ws.Route(ws.DELETE("/orders").To(handler.handleDeleteOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Remove o pedido por completo").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(QueryStruct{}).
		Returns(http.StatusOK, "sucesso", nil).
		Returns(http.StatusBadRequest, "falha", nil))
	ws.Route(ws.POST("/orders/checkout").To(handler.handleCheckout).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Checkout de pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(OrderCheckoutRequest{}).
		Returns(http.StatusOK, "sucesso", Checkout{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	ws.Route(ws.PUT("/orders/status-update").To(handler.handleStatusUpdate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Atualização de status por parte do lojista").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(OrderStatusUpdate{}))
	return handler
}
