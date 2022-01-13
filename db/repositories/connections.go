package repositories

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andriystech/lgc/facilities/ws"
	"github.com/andriystech/lgc/models"
)

var ErrConnIdConflict = errors.New("duplication of connection id")
var ErrConnNotFound = errors.New("connection with provided id not found")

type ConnectionsRepository interface {
	AddConnection(context.Context, string, ws.ConnHelper, *models.User) error
	DeleteConnection(context.Context, string) error
	CountConnections(context.Context) (int, error)
	ConnectedClients(context.Context) ([]string, error)
	GetAllConnections(context.Context) (map[string]ws.ConnHelper, error)
}

type connectionRecord struct {
	conn ws.ConnHelper
	usr  *models.User
}

type connectionsStorage struct {
	db map[string]*connectionRecord
	mu *sync.Mutex
}

func NewConnectionsRepository() ConnectionsRepository {
	return &connectionsStorage{
		db: map[string]*connectionRecord{},
		mu: &sync.Mutex{},
	}
}

func (r *connectionsStorage) AddConnection(ctx context.Context, id string, connection ws.ConnHelper, usr *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db[id] != nil {
		return ErrConnIdConflict
	}
	r.db[id] = &connectionRecord{
		conn: connection,
		usr:  usr,
	}
	return nil
}

func (r *connectionsStorage) DeleteConnection(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db[id] == nil {
		return ErrConnNotFound
	}
	delete(r.db, id)
	return nil
}

func (r *connectionsStorage) CountConnections(ctx context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.db), nil
}

func (r *connectionsStorage) ConnectedClients(ctx context.Context) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res []string = []string{}
	for _, record := range r.db {
		res = append(res, fmt.Sprintf("%s-%s", record.usr.Id, record.usr.UserName))
	}
	return res, nil
}

func (r *connectionsStorage) GetAllConnections(ctx context.Context) (map[string]ws.ConnHelper, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	conns := make(map[string]ws.ConnHelper)
	for _, record := range r.db {
		conns[record.usr.Id] = record.conn
	}
	return conns, nil
}
