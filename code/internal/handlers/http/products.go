package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
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

	ProductList struct {
		Products      []*Product
		Limit, Offset int
		Total         int64
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

	return domain.ParseProductToDomain(
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
	return domain.ParseProductToDomain(
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

	tags := []string{"products"}

	ws.Route(ws.GET("/products/{id}").To(handler.GetProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Obtém dados do produto identificado pelo ID fornecido").
		Param(ws.PathParameter("id", "ID do produto").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(Product{}). // on the response
		Returns(200, "OK", Product{}).
		Returns(500, "ID de produto não cadastrado ou outro erro", nil))

	ws.Route(ws.POST("/products").To(handler.InsertProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionProduct{}). // from the request
		Returns(200, "Produto cadastrado com sucesso", Product{}).
		Returns(500, "Erro ao cadastrar produto", nil))

	ws.Route(ws.PUT("/products").To(handler.UpdateProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Atualiza dados do produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateProduct{}). // from the request
		Returns(200, "Dados do produto atualizados com sucesso", Product{}).
		Returns(500, "Erro ao atualizar dados do produto", nil))

	ws.Route(ws.DELETE("/products").To(handler.DeleteProduct).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Remove produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(DeletionStruct{}). // from the request
		Returns(200, "Produto removido com sucesso", nil).
		Returns(500, "Erro ao remover produto", nil))

	ws.Route(ws.GET("/products").To(handler.ListProductsByCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Lista produtos da categoria especificada").
		Param(ws.QueryParameter("category-id", "ID da categoria").DataType("string")).
		Param(ws.QueryParameter("limit", "Quantidade máxima de entradas que pode retornar").DataType("string")).
		Param(ws.QueryParameter("offset", "Offset a ser usado na paginação").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(Product{}). // on the response
		Returns(200, "OK", Product{}).
		Returns(500, "Erro ao listar produtos", nil))

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
	id := request.QueryParameter("category-id")
	limitS := request.QueryParameter("limit")
	offsetS := request.QueryParameter("offset")

	limit, err := strconv.Atoi(limitS)
	if err != nil {
		return
	}

	offset, err := strconv.Atoi(offsetS)
	if err != nil {
		return
	}

	catId, err := uuid.Parse(id)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	productList, err := pH.productsUC.ListProductsByCategory(pH.ctx, catId, limit, offset)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var prods ProductList
	var prod *Product
	for _, v := range productList.Products {
		prod.fromDomain(v)
		prods.Products = append(prods.Products, prod)
	}
	prods.Total = productList.Total
	prods.Limit = productList.Limit
	prods.Offset = productList.Offset

	_ = response.WriteAsJson(prods)
}
