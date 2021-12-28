package handlers

import (
	"errors"
	"fmt"
	"lgc/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type registerUserHandlerTestData struct {
	payload      string
	wantCode     int
	wantBody     string
	prepareMocks func(*mocks.UserService)
}

type logInUserHandlerTestData struct {
	payload      string
	wantCode     int
	wantBody     string
	prepareMocks func(*mocks.UserService, *mocks.TokenService)
}

type connectionsHandlersTestData struct {
	wantCode     int
	wantBody     string
	prepareMocks func(*mocks.WebSocketService)
}

func TestRegisterUserHandler(t *testing.T) {
	fakeUsr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	ErrCreateNewUsr := errors.New("Unable to create new user")
	ErrSaveNewUsr := errors.New("Unable to save new user")
	testConditions := []registerUserHandlerTestData{
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, fakeUsr.UserName, fakeUsr.Password),
			wantCode: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"id":"%s","userName":"%s"}`, "1", fakeUsr.UserName),
			prepareMocks: func(us *mocks.UserService) {
				us.On("NewUser", fakeUsr.UserName, fakeUsr.Password).Return(fakeUsr, nil)
				us.On("SaveUser", mock.Anything, fakeUsr).Return("1", nil)
			},
		},
		{
			payload:      fmt.Sprintf(`{"userName":"%s,"password":"%s"}`, fakeUsr.UserName, fakeUsr.Password),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"invalid character 'p' after object key:value pair"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService) {},
		},
		{
			payload:      fmt.Sprintf(`{"userName":"%s","password":"%s"}`, fakeUsr.UserName, ""),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'password' was not provided inside body or length less than 6"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService) {},
		},
		{
			payload:      fmt.Sprintf(`{"userName":"%s"}`, fakeUsr.UserName),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'password' was not provided inside body or length less than 6"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService) {},
		},
		{
			payload:      fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "s", fakeUsr.Password),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'userName' was not provided inside body or length less than 3"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService) {},
		},
		{
			payload:      fmt.Sprintf(`{"password":"%s"}`, fakeUsr.Password),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'userName' was not provided inside body or length less than 3"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService) {},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, fakeUsr.UserName, fakeUsr.Password),
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrCreateNewUsr.Error()),
			prepareMocks: func(us *mocks.UserService) {
				us.On("NewUser", fakeUsr.UserName, fakeUsr.Password).Return(nil, ErrCreateNewUsr)
			},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, fakeUsr.UserName, fakeUsr.Password),
			wantCode: http.StatusConflict,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusConflict, repositories.ErrUserWithNameAlreadyExists.Error()),
			prepareMocks: func(us *mocks.UserService) {
				us.On("NewUser", fakeUsr.UserName, fakeUsr.Password).Return(fakeUsr, nil)
				us.On("SaveUser", mock.Anything, fakeUsr).Return("", repositories.ErrUserWithNameAlreadyExists)
			},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, fakeUsr.UserName, fakeUsr.Password),
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrSaveNewUsr.Error()),
			prepareMocks: func(us *mocks.UserService) {
				us.On("NewUser", fakeUsr.UserName, fakeUsr.Password).Return(fakeUsr, nil)
				us.On("SaveUser", mock.Anything, fakeUsr).Return("", ErrSaveNewUsr)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("should respond with %d status and %s body", testCond.wantCode, testCond.wantBody)
		t.Run(tName, func(t *testing.T) {
			us := new(mocks.UserService)
			testCond.prepareMocks(us)
			registerHandler := RegisterUserHandler(us)

			req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(testCond.payload))
			assert.Nil(t, err, "%v", err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCond.wantCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, testCond.wantCode)
			assert.Equal(t, testCond.wantBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.wantBody)
			us.AssertExpectations(t)
		})
	}
}

func TestLogInUserHandler(t *testing.T) {
	fakeToken := "0e903bae-be98-47f3-8d49-e8d950442238"
	fakeUsr := &models.User{UserName: "foobar", Password: "e0b50e3adb85ce07a41196709ba642886ba828a354acee42eae7559bc7c623981897f56020699ed61fa052f4784bf37e76eff016ee065d77bc158dd172eabd76"}
	ErrFindUsrDb := errors.New("Unable to find user")
	ErrTokenGenerate := errors.New("Unable to generate token")
	testConditions := []logInUserHandlerTestData{
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "foobar", "qwerty123456"),
			wantCode: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"url":"ws:///chat/ws.rtm.start?token=%s"}`, fakeToken),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {
				us.On("FindUserByName", mock.Anything, fakeUsr.UserName).Return(fakeUsr, nil)
				ts.On("GenerateToken", mock.Anything, fakeUsr).Return(&models.Token{Payload: fakeToken}, nil)
			},
		},
		{
			payload:      fmt.Sprintf(`{"userName:"%s","password":"%s"}`, "foobar", "qwerty123456"),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"invalid character 'f' after object key"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {},
		},
		{
			payload:      fmt.Sprintf(`{"userName":"%s"}`, "foobar"),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'password' was not provided inside body"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {},
		},
		{
			payload:      fmt.Sprintf(`{"password":"%s"}`, "qwerty123456"),
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"field 'userName' was not provided inside body"}`, http.StatusBadRequest),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "foobar", "qwerty123456"),
			wantCode: http.StatusUnauthorized,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {
				us.On("FindUserByName", mock.Anything, "foobar").Return(nil, repositories.ErrUserNotFound)
			},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "foobar", "e0b50e"),
			wantCode: http.StatusUnauthorized,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {
				usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
				us.On("FindUserByName", mock.Anything, usr.UserName).Return(usr, nil)
			},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "foobar", "qwerty123456"),
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrFindUsrDb.Error()),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {
				us.On("FindUserByName", mock.Anything, "foobar").Return(nil, ErrFindUsrDb)
			},
		},
		{
			payload:  fmt.Sprintf(`{"userName":"%s","password":"%s"}`, "foobar", "qwerty123456"),
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrTokenGenerate.Error()),
			prepareMocks: func(us *mocks.UserService, ts *mocks.TokenService) {
				us.On("FindUserByName", mock.Anything, fakeUsr.UserName).Return(fakeUsr, nil)
				ts.On("GenerateToken", mock.Anything, fakeUsr).Return(nil, ErrTokenGenerate)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("should respond with %d status and %s body", testCond.wantCode, testCond.wantBody)
		t.Run(tName, func(t *testing.T) {
			us := new(mocks.UserService)
			ts := new(mocks.TokenService)
			testCond.prepareMocks(us, ts)
			logInHandler := LogInUserHandler(us, ts)

			req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(testCond.payload))
			assert.Nil(t, err, "%v", err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(logInHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCond.wantCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, testCond.wantCode)
			assert.Equal(t, testCond.wantBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.wantBody)
			us.AssertExpectations(t)
			ts.AssertExpectations(t)
		})
	}
}

func TestActiveConnectionsCountHandler(t *testing.T) {
	ErrConnCount := errors.New("Unable to fetch connections count")
	testConditions := []connectionsHandlersTestData{
		{
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf(`{"count":%d}`, 5),
			prepareMocks: func(wss *mocks.WebSocketService) {
				wss.On("GetActiveConnectionsCount", mock.Anything).Return(5, nil)
			},
		},
		{
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrConnCount.Error()),
			prepareMocks: func(wss *mocks.WebSocketService) {
				wss.On("GetActiveConnectionsCount", mock.Anything).Return(0, ErrConnCount)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("should respond with %d status and %s body", testCond.wantCode, testCond.wantBody)
		t.Run(tName, func(t *testing.T) {
			wssvc := new(mocks.WebSocketService)
			testCond.prepareMocks(wssvc)

			req, err := http.NewRequest(http.MethodGet, "user/active/count", nil)
			assert.Nil(t, err, "%v", err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ActiveConnectionsCountHandler(wssvc))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCond.wantCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, testCond.wantCode)
			assert.Equal(t, testCond.wantBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.wantBody)
			wssvc.AssertExpectations(t)
		})
	}
}

func TestActiveUsersHandler(t *testing.T) {
	users := []string{"1-foobar"}
	ErrActiveUsers := errors.New("Unable to fetch active users")
	testConditions := []connectionsHandlersTestData{
		{
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf(`{"users":["%s"]}`, users[0]),
			prepareMocks: func(wss *mocks.WebSocketService) {
				wss.On("GetActiveUsers", mock.Anything).Return(users, nil)
			},
		},
		{
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, ErrActiveUsers.Error()),
			prepareMocks: func(wss *mocks.WebSocketService) {
				wss.On("GetActiveUsers", mock.Anything).Return(nil, ErrActiveUsers)
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("should respond with %d status and %s body", testCond.wantCode, testCond.wantBody)
		t.Run(tName, func(t *testing.T) {
			wssvc := new(mocks.WebSocketService)
			testCond.prepareMocks(wssvc)

			req, err := http.NewRequest(http.MethodGet, "user/active", nil)
			assert.Nil(t, err, "%v", err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ActiveUsersHandler(wssvc))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCond.wantCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, testCond.wantCode)
			assert.Equal(t, testCond.wantBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.wantBody)
			wssvc.AssertExpectations(t)
		})
	}
}
