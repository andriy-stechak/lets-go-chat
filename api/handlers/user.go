package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andriystech/lgc/db/token"
	userDao "github.com/andriystech/lgc/db/user"
	"github.com/andriystech/lgc/errors"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
)

type RegisterOutput struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	v, httpError := ParseJsonBody(r, &models.User{})
	if httpError != nil {
		httpError.Send(w)
		return
	}
	user := v.(*models.User)
	if httpError := validateUser(user); httpError != nil {
		log.Printf("Invalid input. Reason: %s", httpError.Error())
		httpError.Send(w)
		return
	}
	userId, appError := userDao.Save(user)
	if appError != nil {
		httpError := errors.ToHttpError(appError)
		log.Printf("Unable to save user. Reason: %s", appError.Error())
		httpError.Send(w)
		return
	}
	SendJsonResponse(w, RegisterOutput{Id: userId, UserName: user.UserName}, http.StatusCreated)
}

func LogInUserHandler(w http.ResponseWriter, r *http.Request) {
	c, httpErr := fetchLogInCreds(r)
	if httpErr != nil {
		httpErr.Send(w)
		return
	}
	user, appError := userDao.FindUserByName(c.UserName)
	if appError != nil {
		httpError := errors.ToHttpError(appError)
		errorMessage := fmt.Sprintf("Unable to log in user. Reason: %s", appError.Error())
		log.Print(errorMessage)
		httpError.Send(w)
		return
	}
	if !hasher.CheckPasswordHash(c.Password, user.Password) {
		httpError := errors.HttpUnauthorized("Unable to log in user. Reason: Invalid creds")
		log.Println(httpError.Error())
		httpError.Send(w)
		return
	}
	SendJsonResponse(w, token.Generate(), http.StatusCreated)
}

func fetchLogInCreds(r *http.Request) (*models.UserCreds, *errors.AppHttpError) {
	v, httpError := ParseJsonBody(r, &models.UserCreds{})
	if httpError != nil {
		return nil, httpError
	}
	cred := v.(*models.UserCreds)
	if httpError := validateCreds(cred); httpError != nil {
		log.Printf("Invalid input. Reason: %s", httpError.Error())
		return nil, httpError
	}

	return cred, nil
}

func validateUser(user *models.User) *errors.AppHttpError {
	if len(user.UserName) < models.NameMinLength {
		errorMessage := fmt.Sprintf("Field 'userName' was not provided inside body or length less than %d", models.NameMinLength)
		return errors.HttpBadRequest(errorMessage)
	}
	if len(user.Password) < models.PasswordMinLength {
		errorMessage := fmt.Sprintf("Field 'password' was not provided inside body or length less than %d", models.PasswordMinLength)
		return errors.HttpBadRequest(errorMessage)
	}
	return nil
}

func validateCreds(creds *models.UserCreds) *errors.AppHttpError {
	if len(creds.UserName) == 0 {
		return errors.HttpBadRequest("Field 'userName' was not provided inside body")
	}
	if len(creds.Password) == 0 {
		return errors.HttpBadRequest("Field 'password' was not provided inside body")
	}
	return nil
}
