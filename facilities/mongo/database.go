package mongo

import "go.mongodb.org/mongo-driver/mongo"

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type mongoDatabase struct {
	db *mongo.Database
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}
