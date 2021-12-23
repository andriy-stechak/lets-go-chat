package repositories

import (
	"context"
	"testing"

	"github.com/andriystech/lgc/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestAddConnection(t *testing.T) {
	type testData struct {
		want   error
		connId string
	}

	testConditions := []testData{
		{
			connId: "someid",
			want:   nil,
		},
		{
			connId: "someid",
			want:   ErrConnIdConflict,
		},
		{
			connId: "someid2",
			want:   nil,
		},
	}
	ctx := context.Background()
	repo := NewConnectionsRepository()

	for _, testCond := range testConditions {
		got := repo.AddConnection(ctx, testCond.connId, &websocket.Conn{}, &models.User{})

		assert.Equal(t, testCond.want, got, "AddConnection returned unexpected result: got %v want %v", got, testCond.want)
	}
}

func TestDeleteConnection(t *testing.T) {
	type testData struct {
		want        error
		connIdToAdd string
		connIdToDel string
	}

	testConditions := []testData{
		{
			connIdToAdd: "someid",
			connIdToDel: "someid",
			want:        nil,
		},
		{
			connIdToAdd: "someid2",
			connIdToDel: "someid3",
			want:        ErrConnNotFound,
		},
		{
			connIdToAdd: "someid",
			connIdToDel: "someid2",
			want:        nil,
		},
	}
	ctx := context.Background()
	repo := NewConnectionsRepository()

	for _, testCond := range testConditions {
		repo.AddConnection(ctx, testCond.connIdToAdd, &websocket.Conn{}, &models.User{})

		got := repo.DeleteConnection(ctx, testCond.connIdToDel)

		assert.Equal(t, testCond.want, got, "DeleteConnection returned unexpected result: got %v want %v", got, testCond.want)
	}
}

func TestCountConnectionsSuccess(t *testing.T) {
	ctx := context.Background()
	repo := NewConnectionsRepository()
	repo.AddConnection(ctx, "someid", &websocket.Conn{}, &models.User{})

	gotCount, gotErr := repo.CountConnections(ctx)

	assert.Nil(t, gotErr, "CountConnections returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, 1, gotCount, "CountConnections returned unexpected result: got %v want %v", gotCount, 1)
}

func TestConnectedClientsSuccess(t *testing.T) {
	ctx := context.Background()
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
