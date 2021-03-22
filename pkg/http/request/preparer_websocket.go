package request

import (
	"net/http"
)

// WebSocketHeaders enforce headers to be case-insensitive

// Even if the websocket RFC says that headers should be case-insensitive,
// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
// https://tools.ietf.org/html/rfc6455#page-20
func WebSocketHeaders() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		req.Header["Sec-WebSocket-Key"] = req.Header["Sec-Websocket-Key"]
		req.Header["Sec-WebSocket-Extensions"] = req.Header["Sec-Websocket-Extensions"]
		req.Header["Sec-WebSocket-Accept"] = req.Header["Sec-Websocket-Accept"]
		req.Header["Sec-WebSocket-Protocol"] = req.Header["Sec-Websocket-Protocol"]
		req.Header["Sec-WebSocket-Version"] = req.Header["Sec-Websocket-Version"]
		delete(req.Header, "Sec-Websocket-Key")
		delete(req.Header, "Sec-Websocket-Extensions")
		delete(req.Header, "Sec-Websocket-Accept")
		delete(req.Header, "Sec-Websocket-Protocol")
		delete(req.Header, "Sec-Websocket-Version")
		return req, nil
	})
}
