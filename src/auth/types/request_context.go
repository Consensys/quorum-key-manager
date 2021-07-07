package types

import (
	"crypto/tls"
	"net/http"
)

// UserContext is a set of data attached to every incoming request
type UserContext struct {
	// TLS records information about the TLS connection on which the request was received
	TLS *tls.ConnectionState

	// RemoteAddr records the network address that sent the request
	RemoteAddr string

	// Host records the host on which the URL is sought.
	// This is either the value of the "Host" header or the host name given in the URL itself.
	Host string

	// UserAgent records client's User-Agent
	UserAgent string

	// UserInfo records user information
	// ImpersonatedUserInfo in case the Request impersonate another user
	UserInfo *UserInfo
}

func NewUserContext(req *http.Request) *UserContext {
	return &UserContext{
		TLS:        req.TLS,
		RemoteAddr: req.RemoteAddr,
		Host:       req.Host,
		UserAgent:  req.UserAgent(),
	}
}
