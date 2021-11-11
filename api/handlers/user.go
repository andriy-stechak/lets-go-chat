package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/andriystech/lgc/db/tokens"
	"github.com/andriystech/lgc/db/users"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
)

type RegisterOutput struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	v, err := ParseJsonBody(r, &models.User{})
	if err != nil {
		sendErrorJsonResponse(w, 400, err.Error())
		return
	}
	user := v.(*models.User)
	if err := validateUser(user); err != nil {
		log.Printf("Invalid input. Reason: %s", err.Error())
		sendErrorJsonResponse(w, 400, err.Error())
		return
	}
	userId, err := users.SaveUser(user)
	if err != nil {
		if errors.Is(err, users.ErrUserWithNameAlreadyExists) {
			sendErrorJsonResponse(w, 409, err.Error())
		} else {
			sendErrorJsonResponse(w, 500, err.Error())
		}
		return
	}
	sendJsonResponse(w, RegisterOutput{Id: userId, UserName: user.UserName}, http.StatusCreated)
}

func LogInUserHandler(w http.ResponseWriter, r *http.Request) {
	c, err := fetchLogInCreds(r)
	if err != nil {
		sendErrorJsonResponse(w, 400, err.Error())
		return
	}
	user, err := users.FindUserByName(c.UserName)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			sendErrorJsonResponse(w, 401, "Unable to log in user. Reason: Invalid creds")
		} else {
			sendErrorJsonResponse(w, 500, err.Error())
		}
		return
	}
	if !hasher.CheckPasswordHash(c.Password, user.Password) {
		sendErrorJsonResponse(w, 401, "Unable to log in user. Reason: Invalid creds")
		return
	}
	sendJsonResponse(w, tokens.NewToken(), http.StatusCreated)
}

func fetchLogInCreds(r *http.Request) (*models.UserCreds, error) {
	v, err := ParseJsonBody(r, &models.UserCreds{})
	if err != nil {
		return nil, err
	}
	cred := v.(*models.UserCreds)
	if err := validateCreds(cred); err != nil {
		log.Printf("Invalid input. Reason: %s", err.Error())
		return nil, err
	}

	return cred, nil
}

func validateUser(user *models.User) error {
	if len(user.UserName) < models.NameMinLength {
		return fmt.Errorf("field 'userName' was not provided inside body or length less than %d", models.NameMinLength)
	}
	if len(user.Password) < models.PasswordMinLength {
		return fmt.Errorf("field 'password' was not provided inside body or length less than %d", models.PasswordMinLength)
	}
	return nil
}

func validateCreds(creds *models.UserCreds) error {
	if len(creds.UserName) == 0 {
		return errors.New("field 'userName' was not provided inside body")
	}
	if len(creds.Password) == 0 {
		return errors.New("field 'password' was not provided inside body")
	}
	return nil
}
