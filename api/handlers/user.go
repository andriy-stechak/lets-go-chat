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

type RegisterInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserCredsInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	v, err := ParseJsonBody(r, &RegisterInput{})
	if err != nil {
		sendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	userInputData := v.(*RegisterInput)
	if err := validateUserRegistrationData(userInputData); err != nil {
		log.Printf("Invalid input. Reason: %s", err.Error())
		sendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := models.NewUser(userInputData.UserName, userInputData.Password)
	if err != nil {
		sendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
	}
	userId, err := users.SaveUser(user)
	if errors.Is(err, users.ErrUserWithNameAlreadyExists) {
		sendErrorJsonResponse(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		sendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJsonResponse(w, RegisterOutput{Id: userId, UserName: user.UserName}, http.StatusCreated)
}

func LogInUserHandler(w http.ResponseWriter, r *http.Request) {
	c, err := fetchLogInCreds(r)
	if err != nil {
		sendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := users.FindUserByName(c.UserName)
	if errors.Is(err, users.ErrUserNotFound) {
		sendErrorJsonResponse(w, http.StatusUnauthorized, "Unable to log in user. Reason: Invalid creds")
		return
	}
	if err != nil {
		sendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
	}
	if !hasher.CheckPasswordHash(c.Password, user.Password) {
		sendErrorJsonResponse(w, http.StatusUnauthorized, "Unable to log in user. Reason: Invalid creds")
		return
	}
	sendJsonResponse(w, tokens.Generate(), http.StatusCreated)
}

func fetchLogInCreds(r *http.Request) (*UserCredsInput, error) {
	v, err := ParseJsonBody(r, &UserCredsInput{})
	if err != nil {
		return nil, err
	}
	cred := v.(*UserCredsInput)
	if err := validateCreds(cred); err != nil {
		log.Printf("Invalid input. Reason: %s", err.Error())
		return nil, err
	}

	return cred, nil
}

func validateUserRegistrationData(data *RegisterInput) error {
	if len(data.UserName) < models.NameMinLength {
		return fmt.Errorf("field 'userName' was not provided inside body or length less than %d", models.NameMinLength)
	}
	if len(data.Password) < models.PasswordMinLength {
		return fmt.Errorf("field 'password' was not provided inside body or length less than %d", models.PasswordMinLength)
	}
	return nil
}

func validateCreds(creds *UserCredsInput) error {
	if len(creds.UserName) == 0 {
		return errors.New("field 'userName' was not provided inside body")
	}
	if len(creds.Password) == 0 {
		return errors.New("field 'password' was not provided inside body")
	}
	return nil
}
