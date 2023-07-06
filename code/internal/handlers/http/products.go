package http

import (
	"context"
	"fmt"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type (
	ProductsHttpHandler struct {
		ctx        context.Context
		productsUC ports.ProductsUseCase
	}

	Product struct {
		UpdateProduct
		InsertionProduct
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		DeletedAt time.Time `json:"deleted_at,omitempty"`
	}

	InsertionProduct struct {
		UserID      string `json:"user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
	}

	UpdateProduct struct {
		InsertionProduct
		ID string `json:"id,omitempty"`
	}
)

func (p *Product) fromDomain(product *domain.Product) {
	if p == nil {
		p = &Product{}
	}

	price := fmt.Sprintf("R$ %s", product.Price)

	p.ID = product.ID.String()
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
	p.Name = product.Name
	p.Description = product.Description
	p.CategoryID = product.CategoryID.String()
	p.Price = price
}

func (p *Product) toDomain() *domain.Product {
	if p == nil {
		panic("empty product")
	}

	return domain.ParseToDomain(
		helpers.SafeUUIDFromString(p.ID),
		helpers.SafeUUIDFromString(p.CategoryID),
		p.CreatedAt,
		p.UpdatedAt,
		p.DeletedAt,
		p.Name,
		p.Description,
		p.Price,
	)
}

func (iP *InsertionProduct) toDomain() *domain.Product {
	if iP == nil {
		panic("empty product")
	}

	return domain.NewProduct(uuid.New(), helpers.SafeUUIDFromString(iP.CategoryID), iP.Name, iP.Description, iP.Price)
}

func (uP *UpdateProduct) toDomain() *domain.Product {
	if uP == nil {
		panic("empty product")
	}

	var nilTime time.Time
	return domain.ParseToDomain(
		helpers.SafeUUIDFromString(uP.ID),
		helpers.SafeUUIDFromString(uP.CategoryID),
		nilTime,
		nilTime,
		nilTime,
		uP.Name,
		uP.Description,
		uP.Price,
	)
}

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
	id := request.PathParameter("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	result, err := pH.productsUC.GetProduct(pH.ctx, uid)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var prod Product
	prod.fromDomain(result)
	_ = response.WriteAsJson(prod)

}
func (pH *ProductsHttpHandler) InsertProduct(request *restful.Request, response *restful.Response) {
	var iProd InsertionProduct

	if err := request.ReadEntity(&iProd); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(iProd.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	product, err := pH.productsUC.InsertProduct(pH.ctx, uid, iProd.toDomain())
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var prod Product
	prod.fromDomain(product)
	_ = response.WriteAsJson(prod)
}
func (pH *ProductsHttpHandler) UpdateProduct(request *restful.Request, response *restful.Response) {
	var upProd UpdateProduct

	if err := request.ReadEntity(&upProd); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(upProd.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	product, err := pH.productsUC.UpdateProduct(pH.ctx, uid, upProd.toDomain())
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var prod Product
	prod.fromDomain(product)
	_ = response.WriteAsJson(prod)
}
func (pH *ProductsHttpHandler) DeleteProduct(request *restful.Request, response *restful.Response) {
	var dS DeletionStruct

	if err := request.ReadEntity(&dS); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(dS.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	pID, err := uuid.Parse(dS.ID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	if err = pH.productsUC.DeleteProduct(pH.ctx, uid, pID); err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}
func (pH *ProductsHttpHandler) ListProductsByCategory(request *restful.Request, response *restful.Response) {
	panic("implement me")
}
