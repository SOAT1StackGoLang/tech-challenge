package http

import (
	"context"
	"net/http"

	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

type UserHandler struct {
	ctx          context.Context
	usersUseCase ports.UsersUseCase
}

type NewUser struct {
	Document string `xml:"document" json:"document" description:"CPF do cliente"`
	Name     string `xml:"name" json:"name" description:"Nome do cliente"`
	Email    string `xml:"email" json:"email" description:"Email do cliente"`
}

type IDUser struct {
	ID string `xml:"id" json:"id" description:"ID do cliente no sistema"`
}

type User struct {
	ID       string `xml:"id" json:"id" description:"ID do cliente no sistema"`
	Document string `xml:"document" json:"document" description:"CPF do cliente"`
	Name     string `xml:"name" json:"name" description:"Nome do cliente"`
	Email    string `xml:"email" json:"email" description:"Email do cliente"`
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

	tags := []string{"users"}

	ws.Route(ws.POST("").To(handler.Create).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Cadastra cliente").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(NewUser{}). // from the request
		Returns(200, "Cliente cadastrado", User{}).
		Returns(500, "Erro ao cadastrar cliente", nil))

	ws.Route(ws.GET("/validate/{document}").To(handler.Validate).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).
		Doc("Identifica cliente via CPF. Retorna o ID do cliente no sistema, caso haja cliente cadastrado com esse CPF").
		Param(ws.PathParameter("document", "CPF do cliente").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(IDUser{}). // on the response
		Returns(200, "OK", IDUser{}).
		Returns(500, "CPF não cadastrado ou outro erro", nil))

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
