package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andriystech/lgc/api/errors"
	"github.com/andriystech/lgc/api/handlers/common"
	"github.com/andriystech/lgc/models/creds"
	"github.com/andriystech/lgc/models/token"
	"github.com/andriystech/lgc/models/user"
	"github.com/andriystech/lgc/pkg/hasher"
)

type RegisterOutput struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	v, appError := common.ParseJsonBody(r, &user.Model{})
	if appError != nil {
		appError.Send(w)
		return
	}
	u := v.(*user.Model)
	if appError := u.Validate(); appError != nil {
		log.Printf("Invalid input. Reason: %s", appError.Error())
		appError.Send(w)
		return
	}
	userId, appError := u.Save()
	if appError != nil {
		log.Printf("Unable to save user. Reason: %s", appError.Error())
		appError.Send(w)
		return
	}
	common.SendJsonResponse(w, RegisterOutput{Id: userId, UserName: u.UserName}, http.StatusCreated)
}

func LogInUserHandler(w http.ResponseWriter, r *http.Request) {
	c, appError := fetchLogInCreds(r)
	if appError != nil {
		appError.Send(w)
		return
	}
	user, appError := user.FindUserByName(c.UserName)
	if appError != nil {
		errorMessage := fmt.Sprintf("Unable to log in user. Reason: %s", appError.Error())
		log.Print(errorMessage)
		appError.Send(w)
		return
	}
	if !hasher.CheckPasswordHash(c.Password, user.Password) {
		appError := errors.Unauthorized("Unable to log in user. Reason: Invalid creds")
		log.Println(appError.Error())
		appError.Send(w)
		return
	}
	t := token.Token{}
	t.Generate()
	common.SendJsonResponse(w, t, http.StatusCreated)
}

func fetchLogInCreds(r *http.Request) (*creds.User, *errors.AppError) {
	v, appError := common.ParseJsonBody(r, &creds.User{})
	if appError != nil {
		return nil, appError
	}
	c := v.(*creds.User)
	if appError := c.Validate(); appError != nil {
		log.Printf("Invalid input. Reason: %s", appError.Error())
		return nil, appError
	}

	return c, nil
}
