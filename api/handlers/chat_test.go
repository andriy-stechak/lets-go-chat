package handlers

import (
	"errors"
	"fmt"
	"lgc/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andriystech/lgc/models"
	"github.com/andriystech/lgc/services"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConnectionSuccess(t *testing.T) {
	fakeToken := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	cr := new(mocks.ConnectionsRepository)
	wsvc := services.NewWebSocketService(cr, services.NewUpdater())
	ts := new(mocks.TokenService)
	wsHandler := WSConnectHandler(wsvc, ts)

	ts.On("GetUserByToken", mock.Anything, fakeToken).Return(&models.User{UserName: "foo"}, nil)
	cr.On("AddConnection", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	testServer := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer testServer.Close()

	url := fmt.Sprintf("ws%s?token=%s", strings.TrimPrefix(testServer.URL, "http"), fakeToken)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err, "%v", err)
	defer ws.Close()
	cr.AssertExpectations(t)
	ts.AssertExpectations(t)

	type testDataCond struct {
		input    []byte
		expected string
	}
	testData := []testDataCond{
		{
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			input:    []byte("world"),
			expected: "world",
		},
	}

	for _, testCond := range testData {
		err := ws.WriteMessage(websocket.TextMessage, testCond.input)
		assert.Nil(t, err, "%v", err)

		_, p, err := ws.ReadMessage()
		assert.Nil(t, err, "%v", err)
		got := string(p)
		assert.Equal(t, testCond.expected, got, "bad message got %s want %s", got, testCond.expected)
	}
}

func TestNewConnectionMissingToken(t *testing.T) {
	url := "chat/ws.rtm.start"
	cr := new(mocks.ConnectionsRepository)
	wsvc := services.NewWebSocketService(cr, services.NewUpdater())
	ts := new(mocks.TokenService)
	wsHandler := WSConnectHandler(wsvc, ts)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	expected := fmt.Sprintf(`{"status":%d,"message":"Query parameter 'token' is missing"}`, http.StatusBadRequest)
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	cr.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestNewConnectionInvalidToken(t *testing.T) {
	fakeToken := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	wantErr := errors.New("Invalid token was provided")
	url := fmt.Sprintf("chat/ws.rtm.start?token=%s", fakeToken)
	cr := new(mocks.ConnectionsRepository)
	wsvc := services.NewWebSocketService(cr, services.NewUpdater())
	ts := new(mocks.TokenService)
	ts.On("GetUserByToken", mock.Anything, fakeToken).Return(nil, wantErr)
	wsHandler := WSConnectHandler(wsvc, ts)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusForbidden)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusForbidden, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	cr.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestNewConnectionUnableToOpenWSConnection(t *testing.T) {
	fakeToken := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	wantErr := errors.New("Unable to open web socket connection")
	url := fmt.Sprintf("chat/ws.rtm.start?token=%s", fakeToken)
	wsvc := new(mocks.WebSocketService)
	ts := new(mocks.TokenService)
	ts.On("GetUserByToken", mock.Anything, fakeToken).Return(&models.User{}, nil)
	wsvc.On("NewConnection", mock.Anything, mock.Anything, mock.Anything).Return(wantErr)
	wsHandler := WSConnectHandler(wsvc, ts)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusInternalServerError, wantErr.Error())
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	wsvc.AssertExpectations(t)
	ts.AssertExpectations(t)
}
