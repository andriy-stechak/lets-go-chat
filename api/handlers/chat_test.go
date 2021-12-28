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

type newConnectionSuccessTestData struct {
	input    []byte
	expected string
}

type newConnectionFailTestData struct {
	url          string
	wantCode     int
	wantBody     string
	prepareMocks func(*mocks.ConnectionsRepository, *mocks.TokenService, *mocks.WebSocketService)
}

func TestNewConnectionSuccess(t *testing.T) {
	fakeToken := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	cr := new(mocks.ConnectionsRepository)
	ts := new(mocks.TokenService)
	wsvc := services.NewWebSocketService(cr, services.NewUpdater())
	defer cr.AssertExpectations(t)
	defer ts.AssertExpectations(t)
	wsHandler := WSConnectHandler(wsvc, ts)

	ts.On("GetUserByToken", mock.Anything, fakeToken).Return(&models.User{UserName: "foo"}, nil)
	cr.On("AddConnection", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	testServer := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer testServer.Close()

	url := fmt.Sprintf("ws%s?token=%s", strings.TrimPrefix(testServer.URL, "http"), fakeToken)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err, "%v", err)
	defer ws.Close()

	testConditions := []newConnectionSuccessTestData{
		{
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			input:    []byte("world"),
			expected: "world",
		},
	}

	for _, testCond := range testConditions {
		t.Run(fmt.Sprintf("shoul send %v then recieve %v", string(testCond.input), testCond.expected), func(t *testing.T) {
			err := ws.WriteMessage(websocket.TextMessage, testCond.input)
			assert.Nil(t, err, "%v", err)

			_, p, err := ws.ReadMessage()
			assert.Nil(t, err, "%v", err)
			got := string(p)
			assert.Equal(t, testCond.expected, got, "bad message got %s want %s", got, testCond.expected)
		})
	}
}

func TestNewConnectionFail(t *testing.T) {
	fakeToken := "14ef71b2-5d7c-11ec-a0f3-c46516a4fa45"
	testConditions := []newConnectionFailTestData{
		{
			url:          "chat/ws.rtm.start",
			wantCode:     http.StatusBadRequest,
			wantBody:     fmt.Sprintf(`{"status":%d,"message":"Query parameter 'token' is missing"}`, http.StatusBadRequest),
			prepareMocks: func(cr *mocks.ConnectionsRepository, ts *mocks.TokenService, wsvc *mocks.WebSocketService) {},
		},
		{
			url:      fmt.Sprintf("chat/ws.rtm.start?token=%s", fakeToken),
			wantCode: http.StatusForbidden,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"Invalid token was provided"}`, http.StatusForbidden),
			prepareMocks: func(cr *mocks.ConnectionsRepository, ts *mocks.TokenService, wsvc *mocks.WebSocketService) {
				ts.On("GetUserByToken", mock.Anything, fakeToken).Return(nil, errors.New("Invalid token was provided"))
			},
		},
		{
			url:      fmt.Sprintf("chat/ws.rtm.start?token=%s", fakeToken),
			wantCode: http.StatusInternalServerError,
			wantBody: fmt.Sprintf(`{"status":%d,"message":"Unable to open web socket connection"}`, http.StatusInternalServerError),
			prepareMocks: func(cr *mocks.ConnectionsRepository, ts *mocks.TokenService, wsvc *mocks.WebSocketService) {
				ts.On("GetUserByToken", mock.Anything, fakeToken).Return(&models.User{}, nil)
				wsvc.On("NewConnection", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Unable to open web socket connection"))
			},
		},
	}

	for _, testCond := range testConditions {
		tName := fmt.Sprintf("should respond with %d status and %s body", testCond.wantCode, testCond.wantBody)
		t.Run(tName, func(t *testing.T) {
			cr := new(mocks.ConnectionsRepository)
			ts := new(mocks.TokenService)
			wsvc := new(mocks.WebSocketService)
			testCond.prepareMocks(cr, ts, wsvc)
			wsHandler := WSConnectHandler(wsvc, ts)

			req, err := http.NewRequest(http.MethodGet, testCond.url, nil)
			assert.Nil(t, err, "%v", err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(wsHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCond.wantCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, testCond.wantCode)
			assert.Equal(t, testCond.wantBody, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), testCond.wantBody)
			cr.AssertExpectations(t)
			ts.AssertExpectations(t)
		})
	}
}
