package user

import (
	"fmt"
	"log"

	"github.com/andriystech/lgc/internal/app/errors"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/google/uuid"
)

const NameMinLength = 3

const PasswordMinLength = 6

type Model struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

var storage map[string]*Model

func init() {
	storage = make(map[string]*Model)
}

func FindUserByName(name string) (*Model, *errors.AppError) {
	for _, user := range storage {
		if user.UserName == name {
			return user, nil
		}
	}
	return nil, errors.NotFound(fmt.Sprintf("User with name %s not found", name))
}

func (user *Model) Save() (string, *errors.AppError) {
	if _, err := FindUserByName(user.UserName); err == nil {
		appError := errors.Conflict(fmt.Sprintf("User with name %s already exists", user.UserName))
		log.Printf(appError.Message)
		return "", appError
	}
	id, err := uuid.NewUUID()
	if err != nil {
		appError := errors.InternalError(fmt.Sprintf("Unable to generate user id. Reason %s", err.Error()))
		log.Printf(appError.Message)
		return "", appError
	}

	userId := id.String()
	userPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		appError := errors.InternalError(fmt.Sprintf("Unable to hash user password. Reason %s", err.Error()))
		log.Printf(appError.Message)
		return "", appError
	}
	user.Password = userPassword
	storage[userId] = user
	return userId, nil
}

func (user *Model) Validate() *errors.AppError {
	if len(user.UserName) < NameMinLength {
		errorMessage := fmt.Sprintf("Field 'userName' was not provided inside body or length less than %d", NameMinLength)
		return errors.BadRequest(errorMessage)
	}
	if len(user.Password) < PasswordMinLength {
		errorMessage := fmt.Sprintf("Field 'password' was not provided inside body or length less than %d", PasswordMinLength)
		return errors.BadRequest(errorMessage)
	}
	return nil
}
