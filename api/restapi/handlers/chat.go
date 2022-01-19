package handlers

import (
	"net/http"

	"github.com/andriystech/lgc/api/models"
	"github.com/andriystech/lgc/api/restapi/operations/chat"
	"github.com/andriystech/lgc/services"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type ChatHandler interface {
	Start(chat.WsRTMStartParams) middleware.Responder
	GetActiveUsers(chat.GetActiveUsersParams) middleware.Responder
	GetActiveUsersCount(chat.GetActiveUsersCountParams) middleware.Responder
}

type ChatHandlerContainer struct {
	tokenService     services.TokenService
	webSocketService services.WebSocketService
}

func NewChatHandler(ts services.TokenService, ws services.WebSocketService) ChatHandler {
	return &ChatHandlerContainer{tokenService: ts, webSocketService: ws}
}

func (ch *ChatHandlerContainer) Start(params chat.WsRTMStartParams) middleware.Responder {
	um, err := ch.tokenService.GetUserByToken(params.HTTPRequest.Context(), params.Token)
	if err != nil {
		return chat.NewWsRTMStartBadRequest().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		if err := ch.webSocketService.NewConnection(rw, params.HTTPRequest, um); err != nil {
			chat.NewWsRTMStartInternalServerError().WithPayload(&models.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}).WriteResponse(rw, p)
		}
	})
}

func (ch *ChatHandlerContainer) GetActiveUsers(params chat.GetActiveUsersParams) middleware.Responder {
	activeUsers, err := ch.webSocketService.GetActiveUsers(params.HTTPRequest.Context())
	if err != nil {
		return chat.NewGetActiveUsersInternalServerError().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	return chat.NewGetActiveUsersOK().WithPayload(&models.ActiveUsersResponse{
		Users: activeUsers,
	})
}

func (ch *ChatHandlerContainer) GetActiveUsersCount(params chat.GetActiveUsersCountParams) middleware.Responder {
	connectionsCount, err := ch.webSocketService.GetActiveConnectionsCount(params.HTTPRequest.Context())
	if err != nil {
		return chat.NewGetActiveUsersCountInternalServerError().WithPayload(&models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	count := int64(connectionsCount)
	return chat.NewGetActiveUsersCountOK().WithPayload(&models.ActiveUsersCountResponse{
		Count: &count,
	})
}
