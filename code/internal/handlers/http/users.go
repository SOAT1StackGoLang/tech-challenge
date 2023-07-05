package http

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
)

type UserHandler struct {
	ctx          context.Context
	usersUseCase ports.UsersUseCase
}

type User struct {
	ID       string `json:"id"`
	Document string `json:"document"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u *User) fromDomain(user *domain.User) {
	if u == nil {
		u = &User{}
	}

	u.ID = user.ID.String()
	u.Name = user.Name
	u.Document = user.Document
	u.Email = user.Email
}

func (u *User) toDomain() *domain.User {
	if u == nil {
		u = &User{}
	}

	return &domain.User{
		ID:       helpers.SafeUUIDFromString(u.ID),
		Document: u.Document,
		Name:     u.Name,
		Email:    u.Email,
	}
}

func NewUserHandler(
	ctx context.Context,
	useCase ports.UsersUseCase,
	ws *restful.WebService,
) *UserHandler {
	handler := &UserHandler{
		usersUseCase: useCase,
	}
	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(handler.Create).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/validate/{document}").To(handler.Validate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))

	return handler
}

func (uH *UserHandler) Create(req *restful.Request, resp *restful.Response) {
	var user User

	if err := req.ReadEntity(&user); err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	result, err := uH.usersUseCase.CreateUser(uH.ctx, user.Name, user.Document, user.Email)
	if err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	user.fromDomain(result)
	_ = resp.WriteAsJson(user)
	return
}

func (uH *UserHandler) Validate(req *restful.Request, resp *restful.Response) {
	document := req.PathParameter("document")

	result, err := uH.usersUseCase.ValidateUser(uH.ctx, document)
	if err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	ret := map[string]string{
		"id": result.String(),
	}

	_ = resp.WriteAsJson(ret)
	return
}
