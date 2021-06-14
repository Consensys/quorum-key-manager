package websocket

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/gorilla/websocket"
)

var GoingAway = websocket.FormatCloseMessage(websocket.CloseGoingAway, "")

func Forward(ctx context.Context, from, to *websocket.Conn) error {
	logger := log.FromContext(ctx)

	for {
		typ, msg, err := from.ReadMessage()
		if err != nil {
			logger.WithError(err).Debugf("error reading message")
			return err
		}

		err = to.WriteMessage(typ, msg)
		if err != nil {
			logger.WithError(err).Debugf("error writing message")
		}
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
