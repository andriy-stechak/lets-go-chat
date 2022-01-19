package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

type MultiResultHelper interface {
	All(context.Context, interface{}) error
}

type mongoMultiResult struct {
	mc *mongo.Cursor
}

func (mr *mongoMultiResult) All(ctx context.Context, v interface{}) error {
	return mr.mc.All(ctx, v)
}
