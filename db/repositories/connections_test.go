package repositories

import (
	"context"
	"testing"

	"github.com/andriystech/lgc/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestAddConnectionSuccess(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()
	if got := repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{}); got != nil {
		t.Errorf("AddConnection returned unexpected result: got %v want %v", got, nil)
	}
}

func TestAddConnectionDuplicate(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()
	repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{})

	gotErr := repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{})

	assert.Equal(t, ErrConnIdConflict, gotErr, "AddConnection returned unexpected result: got %v want %v", gotErr, ErrConnIdConflict)
}

func TestDeleteConnectionSuccess(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()
	repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{})

	gotErr := repo.DeleteConnection(ctx, "someid")

	assert.Nil(t, gotErr, "DeleteConnection returned unexpected result: got %v want %v", gotErr, nil)
}

func TestDeleteConnectionNotFound(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()

	gotErr := repo.DeleteConnection(ctx, "someid")

	assert.Equal(t, ErrConnNotFound, gotErr, "DeleteConnection returned unexpected result: got %v want %v", gotErr, ErrConnNotFound)
}

func TestCountConnectionsSuccess(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()
	repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{})

	gotCount, gotErr := repo.CountConnections(ctx)

	assert.Nil(t, gotErr, "CountConnections returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, 1, gotCount, "CountConnections returned unexpected result: got %v want %v", gotCount, 1)
}

func TestConnectedClientsSuccess(t *testing.T) {
	ctx := context.TODO()
	repo := NewConnectionsRepository()
	repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{
		Id:       "someid",
		UserName: "somename",
	})
	want := []string{"someid-somename"}

	got, gotErr := repo.ConnectedClients(ctx)

	assert.Nil(t, gotErr, "ConnectedClients returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, got, want, "ConnectedClients returned unexpected result: [\"someid-somename\"]")
}
