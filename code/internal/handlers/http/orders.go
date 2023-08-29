package http

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
)

type OrdersHttpHandler struct {
	ctx      context.Context
	ordersUC ports.OrdersUseCase
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

	dPs := productsToDomainProducts(insertOrder.ProductsIDs)

	order, err := oH.ordersUC.CreateOrder(oH.ctx, helpers.SafeUUIDFromString(insertOrder.UserID), dPs)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)

		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleAddProductsIntoOrder(request *restful.Request, response *restful.Response) {
	var addRequest UpdateOrder
	if err := request.ReadEntity(&addRequest); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id := helpers.SafeUUIDFromString(addRequest.ID)
	uid := helpers.SafeUUIDFromString(addRequest.InsertionOrder.UserID)

	dPs := productsToDomainProducts(addRequest.InsertionOrder.ProductsIDs)

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

	dPs := productsToDomainProducts(removeReq.ProductsIDs)
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
