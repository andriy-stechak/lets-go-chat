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

func TestRegisterUserHandlerSuccess(t *testing.T) {
	usrId := "1"
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	us.On("NewUser", usr.UserName, usr.Password).Return(usr, nil)
	us.On("SaveUser", mock.Anything, usr).Return(usrId, nil)
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	expectedBody := fmt.Sprintf(`{"id":"%s","userName":"%s"}`, usrId, usr.UserName)
	assert.Equal(t, expectedBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	us.AssertExpectations(t)
}

func TestRegisterUserHandlerInvalidJson(t *testing.T) {
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s,"password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	expectedBody := fmt.Sprintf(`{"status":%d,"message":"invalid character 'p' after object key:value pair"}`, http.StatusBadRequest)
	assert.Equal(t, expectedBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	us.AssertExpectations(t)
}

func TestRegisterUserHandlerInvalidRegData(t *testing.T) {
	type testDataCond struct {
		userName        string
		userPassword    string
		expectedPayload string
	}
	testData := []testDataCond{
		{
			userName:        "foobar",
			userPassword:    "",
			expectedPayload: fmt.Sprintf(`{"status":%d,"message":"field 'password' was not provided inside body or length less than 6"}`, http.StatusBadRequest),
		},
		{
			userName:        "foobar",
			expectedPayload: fmt.Sprintf(`{"status":%d,"message":"field 'password' was not provided inside body or length less than 6"}`, http.StatusBadRequest),
		},
		{
			userName:        "s",
			userPassword:    "qwerty123456",
			expectedPayload: fmt.Sprintf(`{"status":%d,"message":"field 'userName' was not provided inside body or length less than 3"}`, http.StatusBadRequest),
		},
		{
			userPassword:    "qwerty123456",
			expectedPayload: fmt.Sprintf(`{"status":%d,"message":"field 'userName' was not provided inside body or length less than 3"}`, http.StatusBadRequest),
		},
	}
	for _, testCond := range testData {
		us := new(mocks.UserService)
		usr := &models.User{UserName: testCond.userName, Password: testCond.userPassword}
		registerHandler := RegisterUserHandler(us)

		payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, usr.Password)
		req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
		assert.Nil(t, err, "%v", err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, testCond.expectedPayload, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.expectedPayload)
		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		us.AssertExpectations(t)
	}
}

func TestRegisterUserHandlerNewUserFailure(t *testing.T) {
	wantErr := errors.New("Unable to create new user")
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	us.On("NewUser", usr.UserName, usr.Password).Return(nil, wantErr)
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
}

func TestRegisterUserHandlerDuplicate(t *testing.T) {
	wantErr := repositories.ErrUserWithNameAlreadyExists
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	us.On("NewUser", usr.UserName, usr.Password).Return(usr, nil)
	us.On("SaveUser", mock.Anything, usr).Return("", wantErr)
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusConflict)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusConflict, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
}

func TestRegisterUserHandlerUnknown(t *testing.T) {
	wantErr := errors.New("Some error")
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	us.On("NewUser", usr.UserName, usr.Password).Return(usr, nil)
	us.On("SaveUser", mock.Anything, usr).Return("", wantErr)
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
}

func TestLogInUserHandlerSuccess(t *testing.T) {
	fakeToken := "0e903bae-be98-47f3-8d49-e8d950442238"
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar", Password: "e0b50e3adb85ce07a41196709ba642886ba828a354acee42eae7559bc7c623981897f56020699ed61fa052f4784bf37e76eff016ee065d77bc158dd172eabd76"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(usr, nil)
	ts.On("GenerateToken", mock.Anything, usr).Return(&models.Token{Payload: fakeToken}, nil)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	expected := fmt.Sprintf(`{"url":"ws://%s/chat/ws.rtm.start?token=%s"}`, req.Host, fakeToken)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
}

func TestLogInUserHandlerInvalidJson(t *testing.T) {
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName:"%s","password":"%s"}`, "foobar", "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	expected := fmt.Sprintf(`{"status":%d,"message":"invalid character 'f' after object key"}`, http.StatusBadRequest)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestLogInUserHandlerMissingArguments(t *testing.T) {
	type testDataCond struct {
		userName        string
		userPassword    string
		expectedPayload string
	}
	testData := []testDataCond{
		{
			userName:        "foobar",
			expectedPayload: `{"status":400,"message":"field 'password' was not provided inside body"}`,
		},
		{
			userPassword:    "qwerty123456",
			expectedPayload: `{"status":400,"message":"field 'userName' was not provided inside body"}`,
		},
	}

	for _, testCond := range testData {
		us := new(mocks.UserService)
		ts := new(mocks.TokenService)

		logInHandler := LogInUserHandler(us, ts)

		payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, testCond.userName, testCond.userPassword)
		req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
		assert.Nil(t, err, "%v", err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(logInHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		assert.Equal(t, testCond.expectedPayload, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.expectedPayload)
		us.AssertExpectations(t)
		ts.AssertExpectations(t)
	}
}

func TestLogInUserHandlerUserNotFound(t *testing.T) {
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(nil, repositories.ErrUserNotFound)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	expected := fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestLogInUserHandlerInvalidCreds(t *testing.T) {
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar", Password: "e0b50e"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(usr, nil)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	expected := fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestLogInUserHandlerUnableToFindUser(t *testing.T) {
	wantErr := errors.New("Some error")
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(nil, wantErr)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestLogInUserHandlerUnableToGenerateToken(t *testing.T) {
	wantErr := errors.New("Some error")
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar", Password: "e0b50e3adb85ce07a41196709ba642886ba828a354acee42eae7559bc7c623981897f56020699ed61fa052f4784bf37e76eff016ee065d77bc158dd172eabd76"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(usr, nil)
	ts.On("GenerateToken", mock.Anything, usr).Return(nil, wantErr)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	us.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestActiveConnectionsCountHandlerSuccess(t *testing.T) {
	wantCount := 5
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveConnectionsCount", mock.Anything).Return(wantCount, nil)

	req, err := http.NewRequest(http.MethodGet, "user/active/count", nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveConnectionsCountHandler(wssvc))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	expected := fmt.Sprintf(`{"count":%d}`, wantCount)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	wssvc.AssertExpectations(t)
}

func TestActiveConnectionsCountHandlerFail(t *testing.T) {
	wantErr := errors.New("Unable to fetch count")
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveConnectionsCount", mock.Anything).Return(0, wantErr)

	req, err := http.NewRequest(http.MethodGet, "user/active/count", nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveConnectionsCountHandler(wssvc))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	wssvc.AssertExpectations(t)
}

func TestActiveUsersHandlerSuccess(t *testing.T) {
	wantUsers := []string{"1-foobar"}
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveUsers", mock.Anything).Return(wantUsers, nil)

	req, err := http.NewRequest(http.MethodGet, "user/active", nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveUsersHandler(wssvc))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	expected := fmt.Sprintf(`{"users":["%s"]}`, wantUsers[0])
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	wssvc.AssertExpectations(t)
}

func TestActiveUsersHandlerFail(t *testing.T) {
	wantErr := errors.New("Unable to fetch active users")
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveUsers", mock.Anything).Return(nil, wantErr)

	req, err := http.NewRequest(http.MethodGet, "user/active", nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveUsersHandler(wssvc))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	wssvc.AssertExpectations(t)
}
