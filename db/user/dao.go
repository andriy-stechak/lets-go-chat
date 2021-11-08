package user

import (
	"fmt"
	"log"

	"github.com/andriystech/lgc/errors"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/google/uuid"
)

var storage map[string]*models.User = make(map[string]*models.User)

func FindUserByName(name string) (*models.User, *errors.AppError) {
	for _, user := range storage {
		if user.UserName == name {
			return user, nil
		}
	}
	return nil, errors.NotFoundError(fmt.Sprintf("User with name %s not found", name))
}

func Save(user *models.User) (string, *errors.AppError) {
	if _, err := FindUserByName(user.UserName); err == nil {
		appError := errors.ConflictError(fmt.Sprintf("User with name %s already exists", user.UserName))
		log.Printf(appError.Message)
		return "", appError
	}
	id, err := uuid.NewUUID()
	if err != nil {
		appError := errors.GenericError(fmt.Sprintf("Unable to generate user id. Reason %s", err.Error()))
		log.Printf(appError.Message)
		return "", appError
	}

	userId := id.String()
	userPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		appError := errors.GenericError(fmt.Sprintf("Unable to hash user password. Reason %s", err.Error()))
		log.Printf(appError.Message)
		return "", appError
	}
	user.Password = userPassword
	storage[userId] = user
	return userId, nil
}
