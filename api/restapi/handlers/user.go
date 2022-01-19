package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andriystech/lgc/api/models"
	"github.com/andriystech/lgc/api/restapi/operations/user"
	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/andriystech/lgc/services"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

type UserHandler interface {
	Register(user.CreateUserParams) middleware.Responder
	Login(user.LoginUserParams) middleware.Responder
}

type UserHandlerContainer struct {
	userService  services.UserService
	tokenService services.TokenService
}

func NewUserHandler(us services.UserService, ts services.TokenService) UserHandler {
	return &UserHandlerContainer{userService: us, tokenService: ts}
}

func (uh *UserHandlerContainer) Register(params user.CreateUserParams) middleware.Responder {
	um, err := uh.userService.NewUser(*params.Body.UserName, *params.Body.Password)
	if err != nil {
		return user.NewCreateUserInternalServerError().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	userId, err := uh.userService.SaveUser(params.HTTPRequest.Context(), um)
	if errors.Is(err, repositories.ErrUserWithNameAlreadyExists) {
		return user.NewCreateUserConflict().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusConflict,
		})
	}
	if err != nil {
		return user.NewCreateUserInternalServerError().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	return user.NewCreateUserOK().WithPayload(&models.CreateUserResponse{
		ID: strfmt.UUID(userId),
	})
}

func (uh *UserHandlerContainer) Login(params user.LoginUserParams) middleware.Responder {
	um, err := uh.userService.FindUserByName(params.HTTPRequest.Context(), *params.Body.UserName)
	if errors.Is(err, repositories.ErrUserNotFound) {
		return user.NewLoginUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: "Unable to log in user. Reason: Invalid creds",
			Status:  http.StatusBadRequest,
		})
	}
	if !hasher.CheckPasswordHash(*params.Body.Password, um.Password) {
		return user.NewLoginUserBadRequest().WithPayload(&models.ErrorResponse{
			Message: "Unable to log in user. Reason: Invalid creds",
			Status:  http.StatusBadRequest,
		})
	}

	token, err := uh.tokenService.GenerateToken(params.HTTPRequest.Context(), um)
	if err != nil {
		return user.NewLoginUserInternalServerError().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	url := fmt.Sprintf("ws://%s/chat/ws.rtm.start?token=%s", params.HTTPRequest.Host, token.Payload)
	return user.NewLoginUserOK().WithPayload(&models.LoginUserResonse{
		URL: &url,
	})
}
