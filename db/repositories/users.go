package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/andriystech/lgc/config"
	"github.com/andriystech/lgc/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserWithNameAlreadyExists = errors.New("user with provided name already exists")

type UsersRepository struct {
	db *mongo.Collection
}

func NewUsersRepository(client *mongo.Client) *UsersRepository {
	serverConfig := config.GetServerConfig()
	return &UsersRepository{
		db: client.Database(serverConfig.DbName).Collection("users"),
	}
}

func (r *UsersRepository) SaveUser(ctx context.Context, user *models.User) (string, error) {
	if _, err := r.FindUserByName(ctx, user.UserName); err == nil {
		return "", ErrUserWithNameAlreadyExists
	}
	res, err := r.db.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Unable to save user data into database. Reason: %s", err.Error())
		return "", err
	}
	return fmt.Sprintf("%v", res.InsertedID), nil
}

func (r *UsersRepository) FindUserByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.FindOne(
		ctx,
		map[string]string{"userName": name},
	).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		log.Printf("Unable to find user. Reason: %s", err.Error())
		return nil, err
	}
	return &user, nil
}
