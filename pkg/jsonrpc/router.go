package jsonrpc

import (
	"errors"
)

// ErrNotFound is returned when no route match is found.
var ErrNotFound = errors.New("no matching route was found")

type Router struct {
	// DefaultHandler to be used when no route matches.
	defaultHandler Handler

	// Routes to be matched, in order.
	routes []*Route

	// Routes matching a given method name for fast look-up
	methodRoutes map[string][]*Route

	matchers []matcher
}

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{methodRoutes: make(map[string][]*Route)}
}

func (r *Router) Match(msg *RequestMsg, match *RouteMatch) bool {
	methodRoutes, ok := r.methodRoutes[msg.Method]
	if ok {
		for _, route := range methodRoutes {
			if route.Match(msg, match) {
				return true
			}
		}
	}

	for _, route := range r.routes {
		if route.Match(msg, match) {
			return true
		}
	}

	// Closest match for a router (includes sub-routers)
	if r.defaultHandler != nil {
		match.Handler = r.defaultHandler
		return true
	}

	match.Err = ErrNotFound

	return false
}

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route {
	// initialize a route with a copy of the parent router's configuration
	route := &Route{methodRoutes: r.methodRoutes, matchers: make([]matcher, len(r.matchers))}
	r.routes = append(r.routes, route)
	copy(route.matchers, r.matchers)
	return route
}

// Handle registers a new route for a given method
func (r *Router) Handle(method string, handler Handler) *Route {
	return r.NewRoute().Method(method).Handle(handler)
}

// HandleFunc registers a new route for a given method
func (r *Router) HandleFunc(method string, f func(ResponseWriter, *RequestMsg)) *Route {
	return r.NewRoute().Method(method).HandleFunc(f)
}

// Version registers a new route with a matcher for given JSON-RPC version
func (r *Router) DefaultHandler(handler Handler) *Router {
	r.defaultHandler = handler
	return r
}

// Method registers a new route with a matcher for given method
func (r *Router) Method(method string) *Route {
	return r.NewRoute().Method(method)
}

// MethodPrefix registers a new route with a matcher for given method
func (r *Router) MethodPrefix(prefix string) *Route {
	return r.NewRoute().MethodPrefix(prefix)
}

// Version registers a new route with a matcher for given JSON-RPC version
func (r *Router) Version(version string) *Route {
	return r.NewRoute().Version(version)
}

// ServeRPC dispatches the handler registered in the matched route.
func (r *Router) ServeRPC(rw ResponseWriter, msg *RequestMsg) {
	var match RouteMatch
	var handler Handler
	if r.Match(msg, &match) {
		handler = match.Handler
	}

	if handler == nil {
		handler = MethodNotFoundHandler()
	}

	handler.ServeRPC(rw, msg)
}
