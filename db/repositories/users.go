package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/andriystech/lgc/facilities/mongo"
	"github.com/andriystech/lgc/models"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserWithNameAlreadyExists = errors.New("user with provided name already exists")

type UsersRepository interface {
	SaveUser(context.Context, *models.User) (string, error)
	FindUserByName(context.Context, string) (*models.User, error)
	FindUsersNotInIdList(context.Context, []string) ([]*models.User, error)
}

type usersRepository struct {
	db mongo.UsersCollection
}

func NewUsersRepository(db mongo.UsersCollection) UsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) SaveUser(ctx context.Context, user *models.User) (string, error) {
	if _, err := r.FindUserByName(ctx, user.UserName); err == nil {
		return "", ErrUserWithNameAlreadyExists
	}
	res, err := r.db.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Unable to save user data into database. Reason: %s", err.Error())
		return "", err
	}
	return fmt.Sprintf("%v", res), nil
}

func (r *usersRepository) FindUserByName(ctx context.Context, name string) (*models.User, error) {
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

func (r *usersRepository) FindUsersNotInIdList(ctx context.Context, ids []string) ([]*models.User, error) {
	res, err := r.db.Find(
		ctx,
		bson.M{"_id": bson.M{"$nin": ids}},
	)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	if err = res.All(ctx, &users); err != nil {
		if err == mongo.ErrNoDocuments {
			return users, nil
		}
		return nil, err
	}

	return users, nil
}
