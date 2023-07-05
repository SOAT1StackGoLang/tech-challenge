package users

import (
	"context"
	"net/http"

	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

var (
	ctx context.Context
)

type UserHandler struct {
	usersUseCase ports.UsersUseCase
}

func NewUserHandler(
	mainCtx context.Context,
	useCase ports.UsersUseCase,
	ws *restful.WebService,
) *UserHandler {
	ctx = mainCtx
	handler := &UserHandler{
		usersUseCase: useCase,
	}
	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"users"}

	ws.Route(ws.POST("").To(handler.Create).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra cliente").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.GET("/validate/{document}").To(handler.Validate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Identifica cliente via CPF").
		Param(ws.PathParameter("document", "CPF do cliente").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(User{}). // on the response
		Returns(200, "OK", User{}).
		Returns(404, "Not Found", nil))

	return handler
}

func (uH *UserHandler) Create(req *restful.Request, resp *restful.Response) {
	var user User

	if err := req.ReadEntity(&user); err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	result, err := uH.usersUseCase.CreateUser(ctx, user.Name, user.Document, user.Email)
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

	result, err := uH.usersUseCase.ValidateUser(ctx, document)
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
