package repositories

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andriystech/lgc/models"
	"github.com/gorilla/websocket"
)

var ErrConnIdConflict = errors.New("duplication of connection id")
var ErrConnNotFound = errors.New("connection with provided id not found")

type ConnectionRecord struct {
	conn *websocket.Conn
	usr  *models.User
}

type ConnectionsRepository struct {
	db map[string]*ConnectionRecord
	mu *sync.Mutex
}

func NewConnectionsRepository() *ConnectionsRepository {
	return &ConnectionsRepository{
		db: map[string]*ConnectionRecord{},
		mu: &sync.Mutex{},
	}
}

func (r *ConnectionsRepository) AddConnection(ctx context.Context, id string, connection *websocket.Conn, usr *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db[id] != nil {
		return ErrConnIdConflict
	}
	r.db[id] = &ConnectionRecord{
		conn: connection,
		usr:  usr,
	}
	return nil
}

func (r *ConnectionsRepository) DeleteConnection(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db[id] == nil {
		return ErrConnNotFound
	}
	delete(r.db, id)
	return nil
}

func (r *ConnectionsRepository) CountConnections(ctx context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.db), nil
}

func (r *ConnectionsRepository) ConnectedClients(ctx context.Context) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res []string = []string{}
	for _, record := range r.db {
		res = append(res, fmt.Sprintf("%s-%s", record.usr.Id, record.usr.UserName))
	}
	return res, nil
}
