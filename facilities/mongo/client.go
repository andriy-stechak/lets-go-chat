package mongo

import (
	"context"

	"github.com/andriystech/lgc/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessagesCollection CollectionHelper
type UsersCollection CollectionHelper

var ErrNoDocuments = mongo.ErrNoDocuments

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	Find(context.Context, interface{}) (MultiResultHelper, error)
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type MultiResultHelper interface {
	All(context.Context, interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect(context.Context) error
	Disconnect(context.Context) error
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoMultiResult struct {
	mc *mongo.Cursor
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

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
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

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (mr *mongoMultiResult) All(ctx context.Context, v interface{}) error {
	return mr.mc.All(ctx, v)
}

func NewMessagesCollection(
	client ClientHelper,
	config *config.ServerConfig,
) MessagesCollection {
	return client.Database(config.DbName).Collection("messages")
}

func NewUsersCollection(
	client ClientHelper,
	config *config.ServerConfig,
) UsersCollection {
	return client.Database(config.DbName).Collection("users")
}
