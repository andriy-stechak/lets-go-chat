package mongo

import (
	"context"

	"github.com/andriystech/lgc/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect(context.Context) error
	Disconnect(context.Context) error
}

type mongoClient struct {
	cl *mongo.Client
}

func NewClient(cnf *config.ServerConfig) (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(cnf.MongoDbUrl))

	return &mongoClient{cl: c}, err

}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}
