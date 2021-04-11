package jsonrpc

import (
	"strings"
)

// matchers try to match a request.
type matcher interface {
	Match(*Request, *RouteMatch) bool
}

// RouteMatch stores information about a matched route
type RouteMatch struct {
	Route   *Route
	Handler Handler
	Err     error
}

// versionMatcher matches request with a given version
type versionMatcher struct {
	version string
}

func (matcher *versionMatcher) Match(req *Request, _ *RouteMatch) bool {
	return req.Version() == matcher.version
}

// methodMatcher matches request with a given method
type methodMatcher struct {
	method string
}

func (matcher *methodMatcher) Match(req *Request, _ *RouteMatch) bool {
	return req.Method() == matcher.method
}

// methodMatcher matches request with a given prefix
type methodPrefixMatcher struct {
	prefix string
}

func (matcher *methodPrefixMatcher) Match(req *Request, _ *RouteMatch) bool {
	return strings.HasPrefix(req.Method(), matcher.prefix)
}
