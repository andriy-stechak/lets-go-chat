package main

import (
	"context"
	"time"

	"github.com/andriystech/lgc/api/server"
	"github.com/andriystech/lgc/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	serverConfig := config.GetServerConfig()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(serverConfig.DbConnectionTimeoutInSeconds),
	)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(serverConfig.MongoDbUrl))
	if err != nil {
		panic(err)
	}
	defer db.Disconnect(ctx)
	server.Run(db)
}
