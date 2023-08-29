package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
)

type CategoriesHttpHandler struct {
	ctx               context.Context
	categoriesUseCase ports.CategoriesUseCase
}

func NewCategoriesHttpHandler(ctx context.Context, categoriesUseCase ports.CategoriesUseCase, ws *restful.WebService) *CategoriesHttpHandler {
	handler := &CategoriesHttpHandler{
		ctx:               ctx,
		categoriesUseCase: categoriesUseCase,
	}

	tags := []string{"categories"}

	ws.Route(ws.GET("/categories/{id}").To(handler.handleGetCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Obtém informações sobre categoria de produto").
		Param(ws.PathParameter("id", "ID da categoria de produto").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(Category{}). // on the response
		Returns(200, "OK", Category{}).
		Returns(500, "ID de categoria não cadastrada ou outro erro", nil))

	ws.Route(ws.POST("/categories").To(handler.handleInsertCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra categoria de produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionCategory{}). // from the request
		Returns(200, "Categoria cadastrada", Category{}).
		Returns(500, "Erro ao cadastrar categoria", nil))

	ws.Route(ws.DELETE("/categories").To(handler.handleDeleteCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Remove categoria de produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(QueryStruct{}). // from the request
		Returns(200, "Categoria removida", nil).
		Returns(500, "Erro ao remover categoria", nil))

	ws.Route(ws.POST("/categories/all").To(handler.handleListCategories).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Lista categorias").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(ListRequest{}).
		Returns(http.StatusOK, "sucesso", CategoriesList{}).
		Returns(http.StatusInternalServerError, "falha interna do servidor", nil))
	return handler
}

func (cH *CategoriesHttpHandler) handleGetCategory(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	result, err := cH.categoriesUseCase.GetCategory(cH.ctx, uid)
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	var cat Category
	cat.fromDomain(result)
	_ = response.WriteAsJson(cat)
}

func (cH *CategoriesHttpHandler) handleInsertCategory(request *restful.Request, response *restful.Response) {
	var cat Category

	if err := request.ReadEntity(&cat); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(cat.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	domainCat, err := cH.categoriesUseCase.InsertCategory(cH.ctx, uid, cat.toDomain())
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	cat.fromDomain(domainCat)
	_ = response.WriteAsJson(cat)
}

func (cH *CategoriesHttpHandler) handleDeleteCategory(request *restful.Request, response *restful.Response) {
	var dS QueryStruct

	if err := request.ReadEntity(&dS); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(dS.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	cID, err := uuid.Parse(dS.ID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	if err = cH.categoriesUseCase.DeleteCategory(cH.ctx, uid, cID); err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (cH *CategoriesHttpHandler) handleListCategories(request *restful.Request, response *restful.Response) {
	var lR ListRequest
	if err := request.ReadEntity(&lR); err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(lR.UserID)
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	list, err := cH.categoriesUseCase.ListCategories(cH.ctx, uid, lR.Limit, lR.Offset)

	var cL CategoriesList
	var cat Category

	for _, c := range list.Categories {
		cat.fromDomain(c)
		cL.Categories = append(cL.Categories, cat)
	}
	cL.Total = list.Total
	cL.Limit = list.Limit
	cL.Offset = list.Offset

	_ = response.WriteAsJson(cL)
}
