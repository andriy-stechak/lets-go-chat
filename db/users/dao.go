package users

import (
	"errors"

	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/google/uuid"
)

var storage = make(map[string]*models.User)

var ErrUserNotFound = errors.New("user not found")
var ErrUserWithNameAlreadyExists = errors.New("user with provided name already exists")

func FindUserByName(name string) (*models.User, error) {
	for _, user := range storage {
		if user.UserName == name {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func SaveUser(user *models.User) (string, error) {
	if _, err := FindUserByName(user.UserName); err == nil {
		return "", ErrUserWithNameAlreadyExists
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	userId := id.String()
	userPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = userPassword
	storage[userId] = user
	return userId, nil
}
