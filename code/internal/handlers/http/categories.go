package http

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type (
	CategoriesHttpHandler struct {
		ctx               context.Context
		categoriesUseCase ports.CategoriesUseCase
	}

	Category struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		UserID    string    `json:"user_id"`
	}
)

func NewCategoriesHttpHandler(ctx context.Context, categoriesUseCase ports.CategoriesUseCase, ws *restful.WebService) *CategoriesHttpHandler {
	handler := &CategoriesHttpHandler{
		ctx:               ctx,
		categoriesUseCase: categoriesUseCase,
	}

	ws.Route(ws.GET("/categories/{id}").To(handler.GetCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/categories").To(handler.InsertCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.DELETE("/categories").To(handler.DeleteCategory).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	return &CategoriesHttpHandler{ctx: ctx, categoriesUseCase: categoriesUseCase}
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
