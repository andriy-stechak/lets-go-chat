package ws

import "github.com/gorilla/websocket"

type ConnHelper interface {
	Close()
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
}

type websocketConnection struct {
	c *websocket.Conn
}

func NewConn(c *websocket.Conn) ConnHelper {
	return &websocketConnection{
		c: c,
	}
}

func (wc *websocketConnection) Close() {
	wc.c.Close()
}

func (wc *websocketConnection) ReadMessage() (int, []byte, error) {
	return wc.c.ReadMessage()
}

func (wc *websocketConnection) WriteMessage(mt int, msg []byte) error {
	return wc.c.WriteMessage(mt, msg)
}
