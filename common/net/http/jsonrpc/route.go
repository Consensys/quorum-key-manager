package jsonrpc

// Route holds information about a json-rpc route
type Route struct {
	handler Handler

	matchers []Matcher

	// Routes matching a given method name for fast look-up
	methodRoutes map[string][]*Route

	err error
}

func (r *Route) Match(hctx *Context, match *RouteMatch) bool {
	for _, m := range r.matchers {
		if matched := m.Match(hctx, match); !matched {
			return false
		}
	}

	if match.Route == nil {
		match.Route = r
	}

	return true
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m Matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

func (r *Route) Error() error {
	return r.err
}

func (r *Route) Method(method string) *Route {
	r.addMatcher(NewMethodMatcher(method))
	r.methodRoutes[method] = append(r.methodRoutes[method], r)
	return r
}

func (r *Route) MethodPrefix(prefix string) *Route {
	r.addMatcher(NewMethodPrefixMatcher(prefix))
	return r
}

func (r *Route) Version(version string) *Route {
	r.addMatcher(NewVersionMatcher(version))
	return r
}

func (r *Route) Handle(h Handler) *Route {
	r.handler = h
	return r
}

func (r *Route) HandleFunc(f func(*Context)) *Route {
	return r.Handle(HandlerFunc(f))
}

func (r *Route) Subrouter() *Router {
	// initialize a subrouter with a copy of the parent route's configuration
	router := &Router{methodRoutes: r.methodRoutes, matchers: make([]Matcher, len(r.matchers))}
	r.addMatcher(router)
	copy(router.matchers[:], r.matchers)
	return router
}
