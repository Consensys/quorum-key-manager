package integrationtests

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// SignalListener listen to signals and trigger callbacks
type SignalListener struct {
	signals chan os.Signal

	closed    chan struct{}
	closeOnce *sync.Once

	cb func(signal os.Signal)
}

// NewSignalListener creates a new SignalListener
func NewSignalListener(cb func(os.Signal)) *SignalListener {
	l := &SignalListener{
		signals:   make(chan os.Signal, 3),
		closed:    make(chan struct{}),
		closeOnce: &sync.Once{},
		cb:        cb,
	}

	go l.listen()

	return l
}

// Close signal listener
func (l *SignalListener) Close() {
	l.closeOnce.Do(func() {
		close(l.closed)
	})
}

// Listen start Listening to signals
func (l *SignalListener) listen() {
	// Redirect signals
	signal.Notify(l.signals)
signalLoop:
	for {
		select {
		case sig := <-l.signals:
			l.processSignal(sig)
		case <-l.closed:
			break signalLoop
		}
	}
}

func (l *SignalListener) processSignal(sig os.Signal) {
	switch sig {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		log.Warnf("signal: %q intercepted", sig.String())
		l.cb(sig)
	case syscall.SIGPIPE:
		// Ignore random broken pipe
		log.Debugf("signal: %q intercepted", sig.String())
	}
}

func WaitForServiceLive(ctx context.Context, url, name string, timeout time.Duration) {
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		req, _ := http.NewRequest("GET", url, nil)
		req = req.WithContext(rctx)

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			if resp != nil && resp.StatusCode == 200 {
				log.Infof("service %s is live", name)
				return
			}

			log.WithField("status", resp.StatusCode).Warnf("cannot reach %s service", name)
		}

		if rctx.Err() != nil {
			log.WithError(rctx.Err()).Warnf("cannot reach %s service", name)
			return
		}

		log.Debugf("waiting for 1 s for service %s to start...", name)
		time.Sleep(time.Second)
	}
}
