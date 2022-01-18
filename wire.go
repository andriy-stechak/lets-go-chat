//go:build wireinject
// +build wireinject

package main

import (
	"github.com/andriystech/lgc/api/server"
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

func NewServer(db mongo.ClientHelper) server.HttpServer {
	wire.Build(
		config.GetServerConfig,
		ws.NewUpgrader,
		collectionsSet,
		repositoriesSet,
		servicesSet,
		server.NewHttpServer,
	)
	return &server.HttpServerContainer{}
}
