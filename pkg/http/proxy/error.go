package proxy

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
)

type HandleRoundTripErrorFunc func(rw http.ResponseWriter, req *http.Request, err error)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection
const StatusClientClosedRequestText = "Client Closed Connection"

func HandleRoundTripError(rw http.ResponseWriter, req *http.Request, err error) {
	statusCode := StatusCodeFromRoundTripError(err)

	rw.WriteHeader(statusCode)
	_, _ = rw.Write([]byte(StatusText(statusCode)))
}

func StatusCodeFromRoundTripError(err error) int {
	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, io.EOF):
		statusCode = http.StatusBadGateway
	case errors.Is(err, context.Canceled):
		statusCode = StatusClientClosedRequest
	case errors.Is(err, context.DeadlineExceeded):
		statusCode = StatusClientClosedRequest
	default:
		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				statusCode = http.StatusGatewayTimeout
			} else {
				statusCode = http.StatusBadGateway
			}
		}
	}

	return statusCode
}

func StatusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}
