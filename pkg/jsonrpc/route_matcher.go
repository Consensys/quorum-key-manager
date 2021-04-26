package jsonrpc

import (
	"strings"
)

// matchers try to match a msguest.
type matcher interface {
	Match(*RequestMsg, *RouteMatch) bool
}

// RouteMatch stores information about a matched route
type RouteMatch struct {
	Route   *Route
	Handler Handler
	Err     error
}

// versionMatcher matches msguest with a given version
type versionMatcher struct {
	version string
}

func (matcher *versionMatcher) Match(msg *RequestMsg, _ *RouteMatch) bool {
	return msg.Version == matcher.version
}

// methodMatcher matches msguest with a given method
type methodMatcher struct {
	method string
}

func (matcher *methodMatcher) Match(msg *RequestMsg, _ *RouteMatch) bool {
	return msg.Method == matcher.method
}

// methodMatcher matches msguest with a given prefix
type methodPrefixMatcher struct {
	prefix string
}

func (matcher *methodPrefixMatcher) Match(msg *RequestMsg, _ *RouteMatch) bool {
	return strings.HasPrefix(msg.Method, matcher.prefix)
}
