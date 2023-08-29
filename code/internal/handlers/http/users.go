package http

import (
	"context"
	"net/http"

	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

type UserHandler struct {
	ctx          context.Context
	usersUseCase ports.UsersUseCase
}

func NewUserHandler(
	_ context.Context,
	useCase ports.UsersUseCase,
	ws *restful.WebService,
) *UserHandler {
	handler := &UserHandler{
		usersUseCase: useCase,
	}

	tags := []string{"users"}

	ws.Route(ws.POST("/users").To(handler.handleCreate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra cliente").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(InsertionUser{}). // from the request
		Returns(200, "Cliente cadastrado", User{}).
		Returns(500, "Erro ao cadastrar cliente", nil))

	ws.Route(ws.POST("/users/validate").To(handler.handleValidate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Identifica cliente via CPF. Retorna o ID do cliente no sistema, caso haja cliente cadastrado com esse CPF").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(QueryUser{}).
		Writes(ValidatedUser{}). // on the response
		Returns(200, "OK", ValidatedUser{}).
		Returns(500, "CPF não cadastrado ou outro erro", nil))

	return handler
}

func (uH *UserHandler) handleCreate(req *restful.Request, resp *restful.Response) {
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

func (uH *UserHandler) handleValidate(req *restful.Request, resp *restful.Response) {
	var queryUser QueryUser
	if err := req.ReadEntity(&queryUser); err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	result, err := uH.usersUseCase.ValidateUser(uH.ctx, queryUser.Document)
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
