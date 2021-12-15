package main

import (
	"context"
	"time"

	"github.com/andriystech/lgc/api/server"
	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/helpers/mongo"
)

func main() {
	serverConfig := config.GetServerConfig()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(serverConfig.DbConnectionTimeoutInSeconds),
	)
	defer cancel()
	db, err := mongo.NewClient(serverConfig)
	db.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Disconnect(ctx)
	server.Run(db)
}
