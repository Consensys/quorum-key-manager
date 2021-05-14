package header

import "net/http"

// WebSocketHeaders enforce WebSocket headers to be case-sensitive

// Even if the websocket RFC says that headers should be case-insensitive,
// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
// https://tools.ietf.org/html/rfc6455#page-20
func WebSocketHeaders(h http.Header) error {
	h["Sec-WebSocket-Key"] = h["Sec-Websocket-Key"]
	h["Sec-WebSocket-Extensions"] = h["Sec-Websocket-Extensions"]
	h["Sec-WebSocket-Accept"] = h["Sec-Websocket-Accept"]
	h["Sec-WebSocket-Protocol"] = h["Sec-Websocket-Protocol"]
	h["Sec-WebSocket-Version"] = h["Sec-Websocket-Version"]
	delete(h, "Sec-Websocket-Key")
	delete(h, "Sec-Websocket-Extensions")
	delete(h, "Sec-Websocket-Accept")
	delete(h, "Sec-Websocket-Protocol")
	delete(h, "Sec-Websocket-Version")
	return nil
}

func DeleteWebSocketHeaders(h http.Header) error {
	delete(h, "Sec-WebSocket-Key")
	delete(h, "Sec-WebSocket-Extensions")
	delete(h, "Sec-WebSocket-Accept")
	delete(h, "Sec-WebSocket-Protocol")
	delete(h, "Sec-WebSocket-Version")
	return nil
}
