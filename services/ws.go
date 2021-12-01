package services

import (
	"context"
	"log"
	"net/http"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	connections *repositories.ConnectionsRepository
	upgrader    websocket.Upgrader
}

func NewWebSocketService(pool *repositories.ConnectionsRepository) *WebSocketService {
	return &WebSocketService{
		connections: pool,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Disable cross domain check
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (svc *WebSocketService) NewConnection(w http.ResponseWriter, r *http.Request, user *models.User) error {
	c, err := svc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Uable to establish web socket connection. Reason: %s", err.Error())
		return err
	}
	id := uuid.NewString()
	err = svc.connections.AddConnection(r.Context(), id, c, user)
	if err != nil {
		return err
	}
	defer func() {
		c.Close()
		err = svc.connections.DeleteConnection(r.Context(), id)
		if err != nil {
			panic(err)
		}
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	return nil
}

func (svc *WebSocketService) GetActiveConnectionsCount(ctx context.Context) (int, error) {
	return svc.connections.CountConnections(ctx)
}

func (svc *WebSocketService) GetActiveUsers(ctx context.Context) ([]string, error) {
	return svc.connections.ConnectedClients(ctx)
}
