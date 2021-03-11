package jsonrpc

import (
	"errors"
)

// ErrNotFound is returned when no route match is found.
var ErrNotFound = errors.New("no matching route was found")

type Router struct {
	// DefaultHandler to be used when no route matches.
	DefaultHandler Handler

	// Routes to be matched, in order.
	routes []*Route

	// Routes matching a given method name for fast look-up
	methodRoutes map[string][]*Route

	matchers []Matcher
}

func (r *Router) Match(htcx *Context, match *RouteMatch) bool {
	methodRoutes, ok := r.methodRoutes[htcx.Method()]
	if ok {
		for _, route := range methodRoutes {
			if route.Match(htcx, match) {
				return true
			}
		}
	}

	for _, route := range r.routes {
		if route.Match(htcx, match) {
			return true
		}
	}

	match.Err = ErrNotFound

	// Closest match for a router (includes sub-routers)
	if r.DefaultHandler != nil {
		match.Handler = r.DefaultHandler
		return true
	}

	return false
}

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route {
	// initialize a route with a copy of the parent router's configuration
	route := &Route{methodRoutes: r.methodRoutes, matchers: make([]Matcher, len(r.matchers))}
	r.routes = append(r.routes, route)
	copy(route.matchers[:], r.matchers)
	return route
}

// Handle registers a new route for a given method
func (r *Router) Handle(method string, handler Handler) *Route {
	return r.NewRoute().Method(method).Handle(handler)
}

// HandleFunc registers a new route for a given method
func (r *Router) HandleFunc(method string, f func(*Context)) *Route {
	return r.NewRoute().Method(method).HandleFunc(f)
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
