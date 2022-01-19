package handlers

import (
	"github.com/andriystech/lgc/api/restapi/operations/chat"
	"github.com/andriystech/lgc/api/restapi/operations/user"
	"github.com/go-openapi/runtime/middleware"
)

type Handlers interface {
	RegisterUser(user.CreateUserParams) middleware.Responder
	LoginUser(user.LoginUserParams) middleware.Responder
	GetActiveUsers(chat.GetActiveUsersParams) middleware.Responder
	GetActiveUsersCount(chat.GetActiveUsersCountParams) middleware.Responder
	StartChat(chat.WsRTMStartParams) middleware.Responder
}

type HandlersContainer struct {
	user UserHandler
	chat ChatHandler
}

func NewHandlers(uh UserHandler, ch ChatHandler) Handlers {
	return &HandlersContainer{user: uh, chat: ch}
}

func (h *HandlersContainer) RegisterUser(params user.CreateUserParams) middleware.Responder {
	return h.user.Register(params)
}

func (h *HandlersContainer) LoginUser(params user.LoginUserParams) middleware.Responder {
	return h.user.Login(params)
}

func (h *HandlersContainer) GetActiveUsers(params chat.GetActiveUsersParams) middleware.Responder {
	return h.chat.GetActiveUsers(params)
}

func (h *HandlersContainer) GetActiveUsersCount(params chat.GetActiveUsersCountParams) middleware.Responder {
	return h.chat.GetActiveUsersCount(params)
}

func (h *HandlersContainer) StartChat(params chat.WsRTMStartParams) middleware.Responder {
	return h.chat.Start(params)
}
