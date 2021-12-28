package services

import (
	"context"
	"testing"

	"github.com/andriystech/lgc/mocks"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveConnectionsCount(t *testing.T) {
	ctx := context.Background()
	cr := new(mocks.ConnectionsRepository)
	count := 1
	cr.On("CountConnections", ctx).Return(count, nil)
	svc := NewWebSocketService(cr, websocket.Upgrader{})

	gotCount, gotErr := svc.GetActiveConnectionsCount(ctx)

	assert.Nil(t, gotErr, "GetActiveConnectionsCount returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, count, gotCount, "GetActiveConnectionsCount returned unexpected result: got count %v want %v", gotCount, count)

	cr.AssertExpectations(t)
}

func TestGetActiveUsers(t *testing.T) {
	ctx := context.Background()
	cr := new(mocks.ConnectionsRepository)
	clients := []string{"1-user", "2-user2"}
	cr.On("ConnectedClients", ctx).Return(clients, nil)
	svc := NewWebSocketService(cr, websocket.Upgrader{})

	gotClients, gotErr := svc.GetActiveUsers(ctx)

	assert.Nil(t, gotErr, "GetActiveUsers returned unexpected result: got error %v want %v", gotErr, nil)
	assert.Equal(t, clients, gotClients, "GetActiveUsers returned unexpected result: got clients %v want %v", gotClients, clients)

	cr.AssertExpectations(t)
}
