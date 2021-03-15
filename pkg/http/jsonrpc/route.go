package jsonrpc

// Route holds information about a json-rpc route
type Route struct {
	handler Handler

	matchers []matcher

	// Routes matching a given method name for fast look-up
	methodRoutes map[string][]*Route
}

func (r *Route) Match(req *Request, match *RouteMatch) bool {
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			return false
		}
	}

	if match.Route == nil {
		match.Route = r
	}

	if match.Handler == nil {
		match.Handler = r.handler
	}

	return true
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m matcher) *Route {
	r.matchers = append(r.matchers, m)
	return r
}

func (r *Route) Method(method string) *Route {
	r.addMatcher(&methodMatcher{method: method})
	if r.methodRoutes != nil {
		r.methodRoutes[method] = append(r.methodRoutes[method], r)
	}

	return r
}

func (r *Route) MethodPrefix(prefix string) *Route {
	r.addMatcher(&methodPrefixMatcher{prefix: prefix})
	return r
}

func (r *Route) Version(version string) *Route {
	r.addMatcher(&versionMatcher{version: version})
	return r
}

func (r *Route) Handle(h Handler) *Route {
	r.handler = h
	return r
}

func (r *Route) HandleFunc(f func(ResponseWriter, *Request)) *Route {
	return r.Handle(HandlerFunc(f))
}

func (r *Route) Subrouter() *Router {
	// initialize a subrouter with a copy of the parent route's configuration
	router := &Router{methodRoutes: r.methodRoutes, matchers: make([]matcher, len(r.matchers))}
	copy(router.matchers[:], r.matchers)

	r.addMatcher(router)

	return router
}
