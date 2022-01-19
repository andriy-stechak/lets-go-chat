//go:build wireinject
// +build wireinject

package restapi

import (
	"github.com/andriystech/lgc/api/restapi/handlers"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/facilities/ws"
	"github.com/andriystech/lgc/services"
	"github.com/google/wire"
)

var collectionsSet = wire.NewSet(
	mongo.NewMessagesCollection,
	mongo.NewUsersCollection,
)

var repositoriesSet = wire.NewSet(
	repositories.NewConnectionsRepository,
	repositories.NewMessagesRepository,
	repositories.NewTokensRepository,
	repositories.NewUsersRepository,
)

var servicesSet = wire.NewSet(
	services.NewTokenService,
	services.NewUserService,
	services.NewWebSocketService,
)

var handlersSet = wire.NewSet(
	handlers.NewUserHandler,
	handlers.NewChatHandler,
)

func InitializeHandlers(db mongo.ClientHelper) handlers.Handlers {
	wire.Build(
		config.GetServerConfig,
		ws.NewUpgrader,
		collectionsSet,
		repositoriesSet,
		servicesSet,
		handlersSet,
		handlers.NewHandlers,
	)
	return &handlers.HandlersContainer{}
}
