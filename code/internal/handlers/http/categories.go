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
	CategoriesHttpHandler struct {
		ctx               context.Context
		categoriesUseCase ports.CategoriesUseCase
	}

	InsertionCategory struct {
		Name   string `json:"name" description:"Nome da categoria de produto"`
		UserID string `json:"user_id,omitempty" description:"ID do usuario criando categoria"`
	}

	Category struct {
		ID        string    `json:"id" description:"ID da categoria de produto"`
		CreatedAt time.Time `json:"created_at" description:"Epoch time em que categoria foi criada"`
		UpdatedAt time.Time `json:"updated_at" description:"Epoch time em que categoria foi modificada"`
		InsertionCategory
	}
)

func NewCategoriesHttpHandler(ctx context.Context, categoriesUseCase ports.CategoriesUseCase, ws *restful.WebService) *CategoriesHttpHandler {
	handler := &CategoriesHttpHandler{
		ctx:               ctx,
		categoriesUseCase: categoriesUseCase,
	}

	tags := []string{"categories"}

	ws.Route(ws.GET("/categories/{id}").To(handler.GetCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Obtém informações sobre categoria de produto").
		Param(ws.PathParameter("id", "ID da categoria de produto").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(Category{}). // on the response
		Returns(200, "OK", Category{}).
		Returns(500, "ID de categoria não cadastrada ou outro erro", nil))

	ws.Route(ws.POST("/categories").To(handler.InsertCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra categoria de produto").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionCategory{}). // from the request
		Returns(200, "Categoria cadastrada", Category{}).
		Returns(500, "Erro ao cadastrar categoria", nil))

	ws.Route(ws.DELETE("/categories").To(handler.DeleteCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	return handler
}

func (cH *CategoriesHttpHandler) GetCategory(request *restful.Request, response *restful.Response) {
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

func (cH *CategoriesHttpHandler) InsertCategory(request *restful.Request, response *restful.Response) {
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

func (cH *CategoriesHttpHandler) DeleteCategory(request *restful.Request, response *restful.Response) {
	idParam := request.QueryParameters("id")
	userParam := request.QueryParameters("user-id")
	if len(idParam) == 0 || len(userParam) == 0 {
		_ = response.WriteError(http.StatusBadRequest, errors.New("invalid query param"))
		return
	}

	uID, err := uuid.Parse(userParam[0])
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}
	cID, err := uuid.Parse(idParam[0])
	if err != nil {
		_ = response.WriteError(http.StatusBadRequest, err)
		return
	}

	if err = cH.categoriesUseCase.DeleteCategory(cH.ctx, uID, cID); err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (c *Category) toDomain() *domain.Category {
	if c == nil {
		c = &Category{}
	}

	return &domain.Category{
		ID:        helpers.SafeUUIDFromString(c.ID),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}

func (c *Category) fromDomain(cat *domain.Category) {
	if c == nil {
		c = &Category{}
	}

	c.ID = cat.ID.String()
	c.Name = cat.Name
	c.CreatedAt = cat.CreatedAt
	c.UpdatedAt = cat.UpdatedAt
}
