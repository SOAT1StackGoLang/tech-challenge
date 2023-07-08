package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type (
	OrdersHttpHandler struct {
		ctx      context.Context
		ordersUC ports.OrdersUseCase
	}

	Order struct {
		PaymentID string    `json:"payment_id,omitempty" description:"ID do pagamento"`
		CreatedAt time.Time `json:"created_at" description:"Data de criação"`
		UpdatedAt time.Time `json:"updated_at,omitempty" description:"Data de atualização"`
		DeletedAt time.Time `json:"deleted_at,omitempty" description:"Data de deleção"`
		Price     string    `json:"price" description:"Preço do pedido"`
		Status    string    `json:"status" description:"Status do pedido"`
		InsertionOrder
		UpdateOrder
	}

	InsertionOrder struct {
		UserID      string      `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []uuid.UUID `json:"products_ids" description:"ID dos produtos"`
	}

	UpdateOrder struct {
		ID          string      `json:"id" description:"ID do Produto"`
		UserID      string      `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []uuid.UUID `json:"products_ids" description:"ID dos produtos"`
	}

	OrderList struct {
		Orders []Order `json:"orders"`
		Limit  int     `json:"limit"`
		Offset int     `json:"offset"`
		Total  int64   `json:"total"`
	}

	OrderListRequest struct {
		UserID string `json:"user_id"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}
)

func (uO *UpdateOrder) fromDomain(order *domain.Order) {
	if uO == nil {
		uO = &UpdateOrder{}
	}

	uO.ID = order.ID.String()
	uO.UserID = order.UserID.String()
	uO.ProductsIDs = order.ProductsIDs
}

func (iO *InsertionOrder) fromDomain(order *domain.Order) {
	if iO == nil {
		iO = &InsertionOrder{}
	}
	iO.UserID = order.UserID.String()
	iO.ProductsIDs = order.ProductsIDs
}

func (o *Order) fromDomain(order *domain.Order) {
	if o == nil {
		o = &Order{}
	}
	var iO InsertionOrder
	iO.fromDomain(order)
	var uO UpdateOrder
	uO.fromDomain(order)
	o.InsertionOrder = iO
	o.UpdateOrder = uO

	var p string
	if p != "" {
		p = helpers.ParseDecimalToString(order.Price)
		o.Price = p
	}
	o.PaymentID = order.PaymentID.String()
	o.CreatedAt = order.CreatedAt
	o.UpdatedAt = order.UpdatedAt
	o.DeletedAt = order.DeletedAt
	o.Status = order.Status
}

//func (oD *Order) toDomain() *domain.Order {
//	var order *domain.Order
//	id := helpers.SafeUUIDFromString(oD.ID)
//	uID := helpers.SafeUUIDFromString(oD.ID)
//	pID := helpers.SafeUUIDFromString(oD.PaymentID)
//	return domain.ParseToDomainOrder(id, uID, pID, nil, )
//}

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

	order, err := oH.ordersUC.CreateOrder(oH.ctx, helpers.SafeUUIDFromString(insertOrder.UserID), insertOrder.ProductsIDs)
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
	uid := helpers.SafeUUIDFromString(addRequest.UserID)

	order, err := oH.ordersUC.InsertProductsIntoOrder(oH.ctx, uid, id, addRequest.ProductsIDs)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var out Order
	out.fromDomain(order)
	_ = response.WriteAsJson(out)
}

func (oH *OrdersHttpHandler) handleRemoveProductsOfOrder(request *restful.Request, response *restful.Response) {
	var addRequest UpdateOrder
	if err := request.ReadEntity(&addRequest); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	id := helpers.SafeUUIDFromString(addRequest.ID)
	uid := helpers.SafeUUIDFromString(addRequest.UserID)

	order, err := oH.ordersUC.RemoveProductFromOrder(oH.ctx, uid, id, addRequest.ProductsIDs)
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

	ws.Route(ws.GET("/orders/{id}").To(handler.handleGetOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Obtém dados do pedido em função do ID fornecido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(QueryStruct{}).
		Returns(http.StatusOK, "ok", Order{}))
	ws.Route(ws.POST("/orders/all").To(handler.handleListOrders).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Lista pedidos").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(OrderListRequest{}).
		Returns(http.StatusOK, "sucesso", OrderList{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	ws.Route(ws.POST("/orders").To(handler.handleCreateOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionOrder{}).
		Returns(http.StatusOK, "sucesso", Order{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	ws.Route(ws.PUT("/orders/remove").To(handler.handleAddProductsIntoOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Adiciona items ao pedido").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateOrder{}).
		Returns(http.StatusOK, "sucesso", Order{}))
	ws.Route(ws.PUT("/orders/add").To(handler.handleRemoveProductsOfOrder).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
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
	return handler
}
