package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/proxy"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/gorilla/websocket"
)

type Proxy struct {
	Upgrader *websocket.Upgrader
	Dialer   *websocket.Dialer

	Interceptor func(req *http.Request, clientConn, serverConn *websocket.Conn)

	ReqPreparer request.Preparer

	ErrorHandler proxy.HandleRoundTripErrorFunc

	PingPongTimeout        time.Duration
	WriteControlMsgTimeout time.Duration
}

func (prx *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	if websocket.IsWebSocketUpgrade(req) {
		// If websocket we serve websocket
		prx.serveWS(rw, req)
	} else {
		next.ServeHTTP(rw, req)
	}
}

func (prx *Proxy) serveWS(rw http.ResponseWriter, req *http.Request) {
	clientConn, serverConn, err := prx.handleUpgrade(rw, req)
	if err != nil {
		// handleUpgrade has already write response to client so no need to write here
		return
	}

	// Pipe control messaged
	prx.pipeControlMessages(req.Context(), clientConn, serverConn)

	// Triggers Interceptor
	prx.interceptor()(req, clientConn, serverConn)
}

func (prx *Proxy) interceptor() func(req *http.Request, clientConn, serverConn *websocket.Conn) {
	if prx.Interceptor != nil {
		return prx.Interceptor
	}

	return Forward
}

func (prx *Proxy) errorHandler() proxy.HandleRoundTripErrorFunc {
	if prx.ErrorHandler != nil {
		return prx.ErrorHandler
	}

	return proxy.HandleRoundTripError
}

func deleteStandardWebSocketHeader(header http.Header) {
	delete(header, "Sec-WebSocket-Key")
	delete(header, "Sec-WebSocket-Extensions")
	delete(header, "Sec-WebSocket-Protocol")
	delete(header, "Sec-WebSocket-Version")
	delete(header, "Sec-WebSocket-Accept")
}

func (prx *Proxy) handleUpgrade(rw http.ResponseWriter, req *http.Request) (clientConn, serverConn *websocket.Conn, err error) {
	// Prepare request
	outReq := req.Clone(req.Context())
	if prx.ReqPreparer != nil {
		outReq, err = prx.ReqPreparer.Prepare(outReq)
		if err != nil {
			prx.errorHandler()(rw, outReq, err)
			return
		}
	}

	outReq, _ = request.RemoveConnectionHeaders().Prepare(outReq)
	outReq, _ = request.RemoveHopByHopHeaders().Prepare(outReq)
	outReq, _ = request.ForwardedFor().Prepare(outReq)

	// delete headers that will be re-populated on Dial
	deleteStandardWebSocketHeader(outReq.Header)

	outReq.URL.Scheme = "ws"

	// Dial server
	var resp *http.Response
	serverConn, resp, err = prx.Dialer.DialContext(outReq.Context(), outReq.URL.String(), outReq.Header)
	if err != nil {
		prx.errorHandler()(rw, outReq, err)
		return
	}

	// delete headers that will be re-populated on Upgrade
	deleteStandardWebSocketHeader(resp.Header)

	// Upgrade client connection
	clientConn, err = prx.Upgrader.Upgrade(rw, req, resp.Header)
	if err != nil {
		return
	}

	return
}

func (prx *Proxy) pipeControlMessages(ctx context.Context, clientConn, serverConn *websocket.Conn) {
	logger := log.FromContext(ctx)

	_ = clientConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))
	_ = serverConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))

	clientConn.SetPingHandler(func(data string) error {
		logger.WithField("data", data).Trace("received client Ping")

		// We received a message from client so we refresh read timeline
		_ = clientConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))

		// Forward Ping to server
		err := serverConn.WriteControl(websocket.PingMessage, []byte(data), time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Ping message on server connection")
		}

		return nil
	})

	clientConn.SetPongHandler(func(data string) error {
		logger.WithField("data", data).Trace("received client Pong")

		// We received a message from client so we refresh read timeline
		_ = clientConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))

		// Forward Pong to server
		err := serverConn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Pong message on server connection")
		}

		return nil
	})

	serverConn.SetPingHandler(func(data string) error {
		logger.WithField("data", data).Trace("received server Ping")

		// We received a message from server so we refresh read timeline
		_ = serverConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))

		// Forward Ping to client
		err := clientConn.WriteControl(websocket.PingMessage, []byte(data), time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Ping message on cient connection")
		}

		return nil
	})

	serverConn.SetPongHandler(func(data string) error {
		logger.WithField("data", data).Trace("received server Pong")

		// We received a message from server so we refresh read timeline
		_ = serverConn.SetReadDeadline(time.Now().Add(prx.PingPongTimeout))

		// Forward pong to client
		err := clientConn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Pong message on cient connection")
		}

		return nil
	})

	clientConn.SetCloseHandler(func(code int, text string) error {
		logger.WithField("text", text).WithField("code", code).Trace("received client Close")

		fmt.Printf("Handle client close cose=%v text=%v\n", code, text)
		message := websocket.FormatCloseMessage(code, text)
		// We are polite we answer Close back to client
		err := clientConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Close message on client connection")
		}

		// Forward Close to server
		err = serverConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Close message on server connection")
		}

		return nil
	})

	serverConn.SetCloseHandler(func(code int, text string) error {
		logger.WithField("text", text).WithField("code", code).Trace("received server Close")

		message := websocket.FormatCloseMessage(code, text)
		// We are polite we answer Close back to server
		err := serverConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Close message on server connection")
		}

		// Forward Close to client
		err = clientConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(prx.WriteControlMsgTimeout))
		if err != nil {
			logger.WithError(err).Debugf("error writing Close message on client connection")
		}

		return nil
	})
}
