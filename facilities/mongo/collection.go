package mongo

import (
	"context"

	"github.com/andriystech/lgc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionHelper interface {
	Find(context.Context, interface{}) (MultiResultHelper, error)
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
}

type MessagesCollection CollectionHelper

func NewMessagesCollection(client ClientHelper, config *config.ServerConfig) MessagesCollection {
	return client.Database(config.DbName).Collection("messages")
}

type UsersCollection CollectionHelper

func NewUsersCollection(client ClientHelper, config *config.ServerConfig) UsersCollection {
	return client.Database(config.DbName).Collection("users")
}

type mongoCollection struct {
	coll *mongo.Collection
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (MultiResultHelper, error) {
	multiResult, err := mc.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &mongoMultiResult{mc: multiResult}, nil
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}
