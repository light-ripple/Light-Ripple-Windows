package websockets

import (
	"github.com/leavengood/websocket"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	Handler: handler,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

// WebsocketV1Entry upgrades a connection to a websocket.
func WebsocketV1Entry(ctx *fasthttp.RequestCtx) {
	upgrader.UpgradeHandler(ctx)
}
