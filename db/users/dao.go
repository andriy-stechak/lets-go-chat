package users

import (
	"errors"

	"github.com/andriystech/lgc/models"
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
	storage[user.Id] = user
	return user.Id, nil
}
