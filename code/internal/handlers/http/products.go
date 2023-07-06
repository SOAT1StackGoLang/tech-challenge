package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"time"
)

type (
	ProductsHttpHandler struct {
		ctx        context.Context
		productsUC ports.ProductsUseCase
	}

	Product struct {
		ID          uuid.UUID `gorm:"id,primaryKey" json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CategoryID  uuid.UUID `json:"category_id"`
		Price       string    `json:"price"`
	}
)

func NewProductsHttpHandler(ctx context.Context, productsUC ports.ProductsUseCase, ws *restful.WebService) *ProductsHttpHandler {
	handler := &ProductsHttpHandler{
		ctx:        ctx,
		productsUC: productsUC,
	}

	ws.Route(ws.GET("/products/{id}").To(handler.GetProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/products").To(handler.InsertProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.PUT("/products").To(handler.UpdateProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.DELETE("/products").To(handler.DeleteProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/products").To(handler.ListProductsByCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))

	return handler
}

func (pH *ProductsHttpHandler) GetProduct(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
func (pH *ProductsHttpHandler) InsertProduct(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
func (pH *ProductsHttpHandler) UpdateProduct(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
func (pH *ProductsHttpHandler) DeleteProduct(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
func (pH *ProductsHttpHandler) ListProductsByCategory(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
