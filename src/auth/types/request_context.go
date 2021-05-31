package types

import (
	"crypto/tls"
	"net/http"
	"time"
)

// RequestContext is a set of data attached to every incoming request
type RequestContext struct {
	// StartedAt records the date at which the request has been started
	startedAt time.Time

	// TLS records information about the TLS connection on which the request was received
	tls *tls.ConnectionState

	// RemoteAddr records the network address that sent the request
	remoteAddr string

	// Host records the host on which the URL is sought.
	// This is either the value of the "Host" header or the host name given in the URL itself.
	host string

	// UserAgent records client's User-Agent
	userAgent string

	// AuthMode records the mode that succeeded to Authenticate the request ('tls', 'api-key', 'oidc' or '')

	// UserInfo records user information
	// ImpersonatedUserInfo in case the Request impersonate another user
	userInfo, impersonatedUserInfo *UserInfo
}

func NewRequestContext(req *http.Request) *RequestContext {
	return &RequestContext{}
}

func (ctx *RequestContext) StartedAt() time.Time {
	return ctx.startedAt
}

func (ctx *RequestContext) TLS() *tls.ConnectionState {
	return ctx.tls
}

func (ctx *RequestContext) RemoteAddr() string {
	return ctx.remoteAddr
}

func (ctx *RequestContext) UserAgent() string {
	return ctx.userAgent
}

func (ctx *RequestContext) Host() string {
	return ctx.host
}

func (ctx *RequestContext) UserInfo() *UserInfo {
	return ctx.userInfo
}

func (ctx *RequestContext) ImpersonatedUserInfo() *UserInfo {
	return ctx.impersonatedUserInfo
}

// UserInfo are extracted from request credentials by authentication middleware
type UserInfo struct {
	//
	authMode string

	// Username identifies the user
	username string

	// Groups indicates the user's membership to collection of users with specific permissions to access
	groups []string

	// Metadata holds some extra information
	metadata map[string]string
}

func NewUserInfo() *UserInfo {
	return &UserInfo{}
}
