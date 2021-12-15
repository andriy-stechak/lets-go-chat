package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/pkg/hasher"
	"github.com/andriystech/lgc/services"
)

type RegisterOutput struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

type LoginOutput struct {
	Url string `json:"url"`
}

type ActiveConnectionsOutput struct {
	Count int `json:"count"`
}

type ActiveUsersOutput struct {
	Users []string `json:"users"`
}

type RegisterInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserCredsInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func RegisterUserHandler(usvc services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := ParseJsonBody(r, &RegisterInput{})
		if err != nil {
			SendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		userInputData := v.(*RegisterInput)
		if err := validateUserRegistrationData(userInputData); err != nil {
			log.Printf("Invalid input. Reason: %s", err.Error())
			SendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := usvc.NewUser(userInputData.UserName, userInputData.Password)
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		userId, err := usvc.SaveUser(r.Context(), user)
		if errors.Is(err, repositories.ErrUserWithNameAlreadyExists) {
			SendErrorJsonResponse(w, http.StatusConflict, err.Error())
			return
		}
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		sendJsonResponse(w, RegisterOutput{Id: userId, UserName: user.UserName}, http.StatusCreated)
	}
}

func LogInUserHandler(usvc services.UserService, tsvc services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := fetchLogInCreds(r)
		if err != nil {
			SendErrorJsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := usvc.FindUserByName(r.Context(), c.UserName)
		if errors.Is(err, repositories.ErrUserNotFound) {
			SendErrorJsonResponse(w, http.StatusUnauthorized, "Unable to log in user. Reason: Invalid creds")
			return
		}
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !hasher.CheckPasswordHash(c.Password, user.Password) {
			SendErrorJsonResponse(w, http.StatusUnauthorized, "Unable to log in user. Reason: Invalid creds")
			return
		}
		token, err := tsvc.GenerateToken(r.Context(), user)
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		sendJsonResponse(w, composeLoginOutput(r, token), http.StatusCreated)
	}
}

func ActiveConnectionsCountHandler(wssvc services.WebSocketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectionsCount, err := wssvc.GetActiveConnectionsCount(r.Context())
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		sendJsonResponse(w, &ActiveConnectionsOutput{
			Count: connectionsCount,
		}, http.StatusOK)
	}
}

func ActiveUsersHandler(wssvc services.WebSocketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeUsers, err := wssvc.GetActiveUsers(r.Context())
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		sendJsonResponse(w, &ActiveUsersOutput{
			Users: activeUsers,
		}, http.StatusOK)
	}
}

func composeLoginOutput(r *http.Request, token *models.Token) *LoginOutput {
	return &LoginOutput{
		Url: fmt.Sprintf("ws://%s/chat/ws.rtm.start?token=%s", r.Host, token.Payload),
	}
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
