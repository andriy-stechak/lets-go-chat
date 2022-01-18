package ws

import (
	"net/http"

	"github.com/andriystech/lgc/config"
	"github.com/gorilla/websocket"
)

type UpgraderHelper interface {
	Upgrade(http.ResponseWriter, *http.Request) (ConnHelper, error)
}

type websocketUpgrader struct {
	updater *websocket.Upgrader
}

func (wu *websocketUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (ConnHelper, error) {
	conn, err := wu.updater.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return NewConn(conn), nil
}

func NewUpgrader(cg *config.ServerConfig) UpgraderHelper {
	return &websocketUpgrader{
		updater: &websocket.Upgrader{
			ReadBufferSize:  cg.WsReadBuffer,
			WriteBufferSize: cg.WsWriteBuffer,
			// Disable cross domain check
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
