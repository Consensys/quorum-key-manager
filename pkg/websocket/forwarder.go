package websocket

import (
	"context"

	"github.com/gorilla/websocket"
)

var GoingAway = websocket.FormatCloseMessage(websocket.CloseGoingAway, "")

func Forward(_ context.Context, from, to *websocket.Conn) error {
	for {
		typ, msg, err := from.ReadMessage()
		if err != nil {
			return err
		}

		_ = to.WriteMessage(typ, msg)
	}
}

func PipeConn(ctx context.Context, clientConn, serverConn *websocket.Conn) (clientErrors, serverErrors <-chan error) {
	clientErrs := make(chan error, 1)
	serverErrs := make(chan error, 1)

	go func() {
		clientErrs <- Forward(ctx, clientConn, serverConn)
		close(clientErrs)
	}()

	go func() {
		serverErrs <- Forward(ctx, serverConn, clientConn)
		close(serverErrs)
	}()

	return clientErrs, serverErrs
}
