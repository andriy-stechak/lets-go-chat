package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MessagesRepository interface {
	SaveMessage(context.Context, *models.Message) (string, error)
	FindUserMessages(context.Context, string) ([]*models.Message, error)
}

type messagesRepository struct {
	db mongo.CollectionHelper
}

func NewMessagesRepository(db mongo.CollectionHelper) MessagesRepository {
	return &messagesRepository{
		db: db,
	}
}

func (r *messagesRepository) SaveMessage(ctx context.Context, msg *models.Message) (string, error) {
	res, err := r.db.InsertOne(ctx, msg)
	if err != nil {
		log.Printf("Unable to save message data into database. Reason: %s", err.Error())
		return "", err
	}
	return fmt.Sprintf("%v", res), nil
}

func (r *messagesRepository) FindUserMessages(ctx context.Context, id string) ([]*models.Message, error) {
	var messages []*models.Message
	res, err := r.db.Find(
		ctx,
		bson.M{
			"recipientId": id,
		},
	)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = res.All(ctx, &results); err != nil {
		if err == mongo.ErrNoDocuments {
			return messages, nil
		}
		return nil, err
	}

	for _, row := range results {
		var msg models.Message
		data, err := bson.Marshal(row)
		if err != nil {
			return nil, err
		}

		if err = bson.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
