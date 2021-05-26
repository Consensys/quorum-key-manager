package auth

import (
	"crypto/tls"
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
	authMode string

	// UserInfo records user information
	// ImpersonatedUserInfo in case the Request impersonate another user
	userInfo, impersonatedUserInfo *UserInfo

	policyResolver PolicyResolver
}

// UserInfo are extracted from request credentials by authentication middleware
type UserInfo struct {
	// Username identifies the user
	username string

	// Groups indicates the user's membership to collection of users with specific permissions to access
	groups []string

	// Metadata holds some extra information
	metadata map[string]string
}
