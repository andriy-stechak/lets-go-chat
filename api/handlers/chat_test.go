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
	cr.On("DeleteConnection", mock.Anything, mock.Anything).Return(nil)

	testServer := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer testServer.Close()

	url := fmt.Sprintf("ws%s?token=%s", strings.TrimPrefix(testServer.URL, "http"), fakeToken)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	for i := 0; i < 5; i++ {
		if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(p) != "hello" {
			t.Fatalf("bad message")
		}
	}
}

func TestNewConnectionMissingToken(t *testing.T) {
	url := "chat/ws.rtm.start"
	cr := new(mocks.ConnectionsRepository)
	wsvc := services.NewWebSocketService(cr, services.NewUpdater())
	ts := new(mocks.TokenService)
	wsHandler := WSConnectHandler(wsvc, ts)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"Query parameter 'token' is missing"}`, http.StatusBadRequest)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	expected := fmt.Sprintf(`{"status":%d,"message":"%s"}`, http.StatusForbidden, wantErr.Error())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wsHandler)
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
