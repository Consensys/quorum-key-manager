package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/header"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/proxy"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/request"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/response"
	"github.com/consensysquorum/quorum-key-manager/pkg/json"
	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	"github.com/gorilla/websocket"
)

type ProxyConfig struct {
	Upgrader *UpgraderConfig `json:"upgrader,omitempty"`
	Dialer   *DialerConfig   `json:"dialer,omitempty"`

	PingPongTimeout        *json.Duration `json:"pingPongTimeout,omitempty"`
	WriteControlMsgTimeout *json.Duration `json:"writeControlMsgTimeout,omitempty"`
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	if cfg.Upgrader == nil {
		cfg.Upgrader = new(UpgraderConfig)
	}
	cfg.Upgrader.SetDefault()

	if cfg.Dialer == nil {
		cfg.Dialer = new(DialerConfig)
	}
	cfg.Dialer.SetDefault()

	if cfg.PingPongTimeout == nil {
		cfg.PingPongTimeout = &json.Duration{Duration: 60 * time.Second}
	}

	if cfg.WriteControlMsgTimeout == nil {
		cfg.WriteControlMsgTimeout = &json.Duration{Duration: time.Second}
	}

	return cfg
}

type InterceptorFunc func(ctx context.Context, clientConn, serverConn *websocket.Conn) (clientErrors, serverErrors <-chan error)

type Proxy struct {
	Interceptor InterceptorFunc

	ReqPreparer  request.Preparer
	RespModifier response.Modifier

	ErrorHandler proxy.HandleRoundTripErrorFunc

	Upgrader *websocket.Upgrader
	Dialer   *websocket.Dialer

	stopOnce sync.Once
	opsWg    sync.WaitGroup
	ops      chan *operation

	stop chan struct{}
	done chan struct{}

	PingPongTimeout        time.Duration
	WriteControlMsgTimeout time.Duration
	CloseGracePeriod       time.Duration
}

func NewProxy(cfg *ProxyConfig) *Proxy {
	return &Proxy{
		Upgrader:               NewUpgrader(cfg.Upgrader),
		Dialer:                 NewDialer(cfg.Dialer),
		PingPongTimeout:        cfg.PingPongTimeout.Duration,
		WriteControlMsgTimeout: cfg.WriteControlMsgTimeout.Duration,
		CloseGracePeriod:       5 * time.Second,
		stop:                   make(chan struct{}),
		done:                   make(chan struct{}),
		ops:                    make(chan *operation),
	}
}

func (prx *Proxy) Start(context.Context) error {
	go prx.processOps()

	return nil
}

func (prx *Proxy) Done() <-chan struct{} {
	return prx.done
}

// Stop proxy
func (prx *Proxy) Stop(ctx context.Context) error {
	prx.stopOnce.Do(func() {
		close(prx.ops)
		close(prx.stop)
	})

	select {
	case <-prx.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (prx *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	prx.serveWS(rw, req)
}

func (prx *Proxy) processOps() {
	for op := range prx.ops {
		prx.opsWg.Add(1)
		go func(op *operation) {
			op.run()
			prx.opsWg.Done()
		}(op)
	}

	prx.opsWg.Wait()
	close(prx.done)
}

func (prx *Proxy) RegisterServerShutdown(srv *http.Server) {
	srv.RegisterOnShutdown(func() {
		_ = prx.Stop(context.Background())
	})
}

func (prx *Proxy) serveWS(rw http.ResponseWriter, req *http.Request) {
	// Upgrade to websocket on client and server conn
	clientConn, serverConn, err := prx.handleUpgrade(rw, req)
	if err != nil {
		return
	}

	// Create a new operation
	op := &operation{
		prx:                 prx,
		req:                 req,
		logger:              log_old.FromContext(req.Context()),
		clientConn:          clientConn,
		serverConn:          serverConn,
		done:                make(chan struct{}),
		receivedClientClose: make(chan struct{}),
		receivedServerClose: make(chan struct{}),
	}

	// sed ops for processing
	prx.ops <- op

	// wait for operation to complete (it avoids request context to be canceled)
	<-op.done
}

func (prx *Proxy) interceptor() InterceptorFunc {
	if prx.Interceptor != nil {
		return prx.Interceptor
	}

	return PipeConn
}

func (prx *Proxy) errorHandler() proxy.HandleRoundTripErrorFunc {
	if prx.ErrorHandler != nil {
		return prx.ErrorHandler
	}

	return proxy.HandleRoundTripError
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

	// Prepare headers for proxying
	outReq, _ = request.RemoveConnectionHeaders().Prepare(outReq)
	outReq, _ = request.RemoveHopByHopHeaders().Prepare(outReq)
	outReq, _ = request.ForwardedFor().Prepare(outReq)

	// delete websocket headers that will be re-populated on Dial
	outReq, _ = request.HeadersPreparer(header.WebSocketHeaders).Prepare(outReq)
	outReq, _ = request.HeadersPreparer(header.DeleteWebSocketHeaders).Prepare(outReq)

	outReq.URL.Scheme = "ws"

	// Dial server
	var resp *http.Response
	serverConn, resp, err = prx.Dialer.DialContext(outReq.Context(), outReq.URL.String(), outReq.Header)
	if err != nil {
		prx.errorHandler()(rw, outReq, err)
		return
	}

	if prx.RespModifier != nil {
		err = prx.RespModifier.Modify(resp)
		if err != nil {
			prx.errorHandler()(rw, outReq, err)
			return
		}
	}

	// delete headers that will be re-populated on Upgrade
	_ = response.HeadersModifier(header.DeleteWebSocketHeaders).Modify(resp)

	// Upgrade client connection
	clientConn, err = prx.Upgrader.Upgrade(rw, req, resp.Header)
	if err != nil {
		_ = prx.writeClose(serverConn, GoingAway, true)
		serverConn.Close()
		return
	}

	return
}

func (prx *Proxy) writeClose(conn *websocket.Conn, msg []byte, wait bool) error {
	// Write close message on connection
	err := prx.writeControl(conn, websocket.CloseMessage, msg)
	if err != nil {
		return err
	}

	if wait {
		// Give a grace period for remote to send CloseMessage back
		time.Sleep(prx.CloseGracePeriod)
	}

	return nil
}

func (prx *Proxy) writeControl(conn *websocket.Conn, typ int, msg []byte) error {
	return conn.WriteControl(typ, msg, time.Now().Add(prx.WriteControlMsgTimeout))
}

type operation struct {
	prx *Proxy

	req    *http.Request
	logger *log_old.Logger

	clientConn, serverConn *websocket.Conn

	writeClientCloseOnce, writeServerCloseOnce sync.Once
	receivedClientClose, receivedServerClose   chan struct{}

	done chan struct{}
}

func (op *operation) run() {
	// Set read timeouts
	_ = op.clientConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))
	_ = op.serverConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))

	// Pipe control messages
	op.pipeControlMessages()

	// Handle stop
	go op.handleStop()

	// Run interceptor
	clientErrs, serverErrs := op.prx.interceptor()(op.req.Context(), op.clientConn, op.serverConn)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err := <-clientErrs
		msg := GoingAway
		if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code != websocket.CloseAbnormalClosure {
			msg = websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
		}
		op.writeServerClose(msg, true)
		wg.Done()
	}()

	go func() {
		err := <-serverErrs
		msg := GoingAway
		if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code != websocket.CloseAbnormalClosure {
			msg = websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
		}
		op.writeClientClose(msg, true)
		wg.Done()
	}()

	wg.Wait()
	close(op.done)
}

func (op *operation) handleStop() {
	<-op.prx.stop
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		op.writeServerClose(GoingAway, true)
		wg.Done()
	}()

	go func() {
		op.writeClientClose(GoingAway, true)
		wg.Done()
	}()

	wg.Wait()
}

func (op *operation) writeClientClose(msg []byte, waitCloseBack bool) {
	op.writeClientCloseOnce.Do(func() {
		err := op.prx.writeControl(op.clientConn, websocket.CloseMessage, msg)
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Close client connection")
			return
		}

		if waitCloseBack {
			after := time.After(op.prx.CloseGracePeriod)
			select {
			case <-op.receivedClientClose:
			case <-after:
			}
		}
	})
}

func (op *operation) writeServerClose(msg []byte, waitCloseBack bool) {
	op.writeServerCloseOnce.Do(func() {
		err := op.prx.writeControl(op.serverConn, websocket.CloseMessage, msg)
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Close to server connection")
			return
		}

		if waitCloseBack {
			after := time.After(op.prx.CloseGracePeriod)
			select {
			case <-op.receivedServerClose:
			case <-after:
			}
		}
	})
}

func (op *operation) pipeControlMessages() {
	op.clientConn.SetPingHandler(func(data string) error {
		op.logger.WithField("data", data).Trace("received client Ping")

		// We received a message from client so we refresh read timeline
		_ = op.clientConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))

		// Forward Ping to server
		err := op.prx.writeControl(op.serverConn, websocket.PingMessage, []byte(data))
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Ping message on server connection")
		}

		return nil
	})

	op.clientConn.SetPongHandler(func(data string) error {
		op.logger.WithField("data", data).Trace("received client Pong")

		// We received a message from client so we refresh read timeline
		_ = op.clientConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))

		// Forward Pong to server
		err := op.prx.writeControl(op.serverConn, websocket.PongMessage, []byte(data))
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Pong message on server connection")
		}

		return nil
	})

	op.serverConn.SetPingHandler(func(data string) error {
		op.logger.WithField("data", data).Trace("received server Ping")

		// We received a message from server so we refresh read timeline
		_ = op.serverConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))

		// Forward Ping to client
		err := op.prx.writeControl(op.clientConn, websocket.PingMessage, []byte(data))
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Ping message on cient connection")
		}

		return nil
	})

	op.serverConn.SetPongHandler(func(data string) error {
		op.logger.WithField("data", data).Trace("received server Pong")

		// We received a message from server so we refresh read timeline
		_ = op.serverConn.SetReadDeadline(time.Now().Add(op.prx.PingPongTimeout))

		// Forward pong to client
		err := op.prx.writeControl(op.clientConn, websocket.PongMessage, []byte(data))
		if err != nil {
			op.logger.WithError(err).Debugf("error writing Pong message on cient connection")
		}

		return nil
	})

	op.clientConn.SetCloseHandler(func(code int, text string) error {
		op.logger.WithField("text", text).WithField("code", code).Trace("received client Close")

		select {
		case <-op.receivedClientClose:
			return nil
		default:
			close(op.receivedClientClose)
		}

		// We answer Close back to client
		msg := []byte{}
		if code != websocket.CloseNoStatusReceived {
			msg = websocket.FormatCloseMessage(code, "")
		}
		op.writeClientClose(msg, false)

		return nil
	})

	op.serverConn.SetCloseHandler(func(code int, text string) error {
		op.logger.WithField("text", text).WithField("code", code).Trace("received server Close")

		select {
		case <-op.receivedServerClose:
			return nil
		default:
			close(op.receivedServerClose)
		}

		// We answer Close back to server
		msg := []byte{}
		if code != websocket.CloseNoStatusReceived {
			msg = websocket.FormatCloseMessage(code, "")
		}
		op.writeServerClose(msg, false)

		return nil
	})
}
