package http

import "net/http"

// Modifier is the interface that allows to modify http.http.Response the Respond method.
type Modifier interface {
	Modify(*http.Response) error
}

// ModifierFunc is a method that implements the Modifier interface.
type ModifierFunc func(*http.Response) error

// Respond implements the Modifier interface on ModifierFunc.
func (rf ModifierFunc) Modify(r *http.Response) error {
	return rf(r)
}

// BackendServer set "X-Backend-Server" header to the the URL of the request
func BackendServer(rest *http.Response) error {
	resp.Header.Set("X-Backend-Server", resp.Request.URL.String())
	return nil
}



// RespondDecorator takes and possibly decorates, by wrapping, a Modifier. Decorators may react to
// the http.http.Response and pass it along or, first, pass the http.http.Response along then react.
type RespondDecorator func(Modifier) Modifier
