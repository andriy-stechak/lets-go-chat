package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/andriystech/lgc/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type connTestData struct {
	want        error
	connId      string
	prepareRepo func(ConnectionsRepository) ConnectionsRepository
}

func TestAddConnection(t *testing.T) {
	testConditions := []connTestData{
		{
			connId: "someid",
			want:   nil,
			prepareRepo: func(cr ConnectionsRepository) ConnectionsRepository {
				return cr
			},
		},
		{
			connId: "someid2",
			want:   ErrConnIdConflict,
			prepareRepo: func(cr ConnectionsRepository) ConnectionsRepository {
				cr.AddConnection(context.Background(), "someid2", &websocket.Conn{}, &models.User{})
				return cr
			},
		},
	}

	for _, testCond := range testConditions {
		t.Run(fmt.Sprintf("AddConnection(%v, %v) == %v", context.Background(), testCond.connId, testCond.want), func(t *testing.T) {
			ctx := context.Background()
			repo := testCond.prepareRepo(NewConnectionsRepository())
			got := repo.AddConnection(ctx, testCond.connId, &websocket.Conn{}, &models.User{})

			assert.Equal(t, testCond.want, got, "AddConnection returned unexpected result: got %v want %v", got, testCond.want)
		})
	}
}

func TestDeleteConnection(t *testing.T) {
	testConditions := []connTestData{
		{
			connId: "someid",
			want:   ErrConnNotFound,
			prepareRepo: func(cr ConnectionsRepository) ConnectionsRepository {
				return cr
			},
		},
		{
			connId: "someid2",
			want:   nil,
			prepareRepo: func(cr ConnectionsRepository) ConnectionsRepository {
				cr.AddConnection(context.Background(), "someid2", &websocket.Conn{}, &models.User{})
				return cr
			},
		},
	}

	for _, testCond := range testConditions {
		t.Run(fmt.Sprintf("DeleteConnection(%v) == %v", testCond.connId, testCond.want), func(t *testing.T) {
			ctx := context.Background()
			repo := testCond.prepareRepo(NewConnectionsRepository())
			got := repo.DeleteConnection(ctx, testCond.connId)

			assert.Equal(t, testCond.want, got, "DeleteConnection returned unexpected result: got %v want %v", got, testCond.want)
		})
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
