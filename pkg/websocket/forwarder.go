package websocket

import (
	"net/http"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/gorilla/websocket"
)

var goingAway = websocket.FormatCloseMessage(websocket.CloseGoingAway, "")

func Forward(req *http.Request, clientConn, serverConn *websocket.Conn) {
	logger := log.FromContext(req.Context())
	go func() {
		defer clientConn.Close()
		for {
			typ, msg, err := clientConn.ReadMessage()
			if err != nil {
				logger.WithError(err).Debugf("error reading message on client connection")

				err = serverConn.WriteControl(websocket.CloseMessage, goingAway, time.Now().Add(time.Second))
				if err != nil {
					logger.WithError(err).Debugf("error writing Close control message on server connection")
				}

				return
			}

			err = serverConn.WriteMessage(typ, msg)
			if err != nil {
				logger.WithError(err).Debugf("error writing message on server connection")
			}
		}
	}()

	go func() {
		defer serverConn.Close()
		for {
			typ, msg, err := serverConn.ReadMessage()
			if err != nil {
				logger.WithError(err).Debugf("error reading message on server connection")

				err = clientConn.WriteControl(websocket.CloseMessage, goingAway, time.Now().Add(time.Second))
				if err != nil {
					logger.WithError(err).Debugf("error writing Close control message on client connection")
				}

				return
			}

			err = clientConn.WriteMessage(typ, msg)
			if err != nil {
				logger.WithError(err).Debugf("error writing message on client connection")
			}
		}
	}()
}
