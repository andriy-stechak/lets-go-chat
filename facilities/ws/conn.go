package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnHelper interface {
	Close()
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
}

type websocketConnection struct {
	c  *websocket.Conn
	mu *sync.Mutex
}

func NewConn(c *websocket.Conn) ConnHelper {
	return &websocketConnection{
		c:  c,
		mu: &sync.Mutex{},
	}
}

func (wc *websocketConnection) Close() {
	wc.c.Close()
}

func (wc *websocketConnection) ReadMessage() (int, []byte, error) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.c.ReadMessage()
}

func (wc *websocketConnection) WriteMessage(mt int, msg []byte) error {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.c.WriteMessage(mt, msg)
}
