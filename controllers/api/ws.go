package api

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var _upg = websocket.Upgrader{}

func init() {
	_upg.CheckOrigin = func(r *http.Request) bool {
		// allow all connections by default
		return true
	}
}

type wsController struct {
}

func (w *wsController) WebSocket(c apiController) (*websocket.Conn, error) {
	return _upg.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
}
