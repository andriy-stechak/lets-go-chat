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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := fmt.Sprintf(`{"id":"%s","userName":"%s"}`, usrId, usr.UserName)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRegisterUserHandlerInvalidJson(t *testing.T) {
	us := new(mocks.UserService)
	usr := &models.User{UserName: "foobar", Password: "qwerty123456"}
	registerHandler := RegisterUserHandler(us)

	payload := fmt.Sprintf(`{"userName":"%s,"password":"%s"}`, usr.UserName, usr.Password)
	req, err := http.NewRequest(http.MethodPost, "user", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"invalid character 'p' after object key:value pair"}`, http.StatusBadRequest)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		if rr.Body.String() != testCond.expectedPayload {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), testCond.expectedPayload)
		}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusConflict, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := fmt.Sprintf(`{"url":"ws://%s/chat/ws.rtm.start?token=%s"}`, req.Host, fakeToken)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLogInUserHandlerInvalidJson(t *testing.T) {
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName:"%s","password":"%s"}`, "foobar", "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"invalid character 'f' after object key"}`, http.StatusBadRequest)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(logInHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		if rr.Body.String() != testCond.expectedPayload {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), testCond.expectedPayload)
		}
	}
}

func TestLogInUserHandlerUserNotFound(t *testing.T) {
	fakeToken := "0e903bae-be98-47f3-8d49-e8d950442238"
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(nil, repositories.ErrUserNotFound)
	ts.On("GenerateToken", mock.Anything, usr).Return(&models.Token{Payload: fakeToken}, nil)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLogInUserHandlerInvalidCreds(t *testing.T) {
	fakeToken := "0e903bae-be98-47f3-8d49-e8d950442238"
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar", Password: "e0b50e"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(usr, nil)
	ts.On("GenerateToken", mock.Anything, usr).Return(&models.Token{Payload: fakeToken}, nil)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"Unable to log in user. Reason: Invalid creds"}`, http.StatusUnauthorized)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLogInUserHandlerUnableToFindUser(t *testing.T) {
	wantErr := errors.New("Some error")
	fakeToken := "0e903bae-be98-47f3-8d49-e8d950442238"
	us := new(mocks.UserService)
	ts := new(mocks.TokenService)
	usr := &models.User{UserName: "foobar"}
	us.On("FindUserByName", mock.Anything, usr.UserName).Return(nil, wantErr)
	ts.On("GenerateToken", mock.Anything, usr).Return(&models.Token{Payload: fakeToken}, nil)

	logInHandler := LogInUserHandler(us, ts)

	payload := fmt.Sprintf(`{"userName":"%s","password":"%s"}`, usr.UserName, "qwerty123456")
	req, err := http.NewRequest(http.MethodPost, "user/login", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logInHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestActiveConnectionsCountHandlerSuccess(t *testing.T) {
	wantCount := 5
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveConnectionsCount", mock.Anything).Return(wantCount, nil)

	req, err := http.NewRequest(http.MethodGet, "user/active/count", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveConnectionsCountHandler(wssvc))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"count":%d}`, wantCount)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestActiveConnectionsCountHandlerFail(t *testing.T) {
	wantErr := errors.New("Unable to fetch count")
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveConnectionsCount", mock.Anything).Return(0, wantErr)

	req, err := http.NewRequest(http.MethodGet, "user/active/count", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveConnectionsCountHandler(wssvc))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestActiveUsersHandlerSuccess(t *testing.T) {
	wantUsers := []string{"1-foobar"}
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveUsers", mock.Anything).Return(wantUsers, nil)

	req, err := http.NewRequest(http.MethodGet, "user/active", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveUsersHandler(wssvc))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"users":["%s"]}`, wantUsers[0])
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestActiveUsersHandlerFail(t *testing.T) {
	wantErr := errors.New("Unable to fetch active users")
	wssvc := new(mocks.WebSocketService)
	wssvc.On("GetActiveUsers", mock.Anything).Return(nil, wantErr)

	req, err := http.NewRequest(http.MethodGet, "user/active", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ActiveUsersHandler(wssvc))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
