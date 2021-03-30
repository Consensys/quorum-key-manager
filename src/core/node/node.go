package node

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/jsonrpc"
)

type Info struct {
	Name string
	Addr string
}

// Node allow to interface with a downstream node
type Node interface {
	// Addr of the downstream node
	Info() *Info

	// Wrap a request and returns a new node from it
	Wrap(req *jsonrpc.Request) Node

	// Do sends an HTTP request and returns an HTTP response, following
	// policy (such as redirects, cookies, auth) as configured for the downstream node
	Do(req *http.Request) (*http.Response, error)

	// RoundTrip executes a single HTTP transaction, returning
	// a Response for the provided Request.
	RoundTrip(*http.Request) (*http.Response, error)

	// Proxy returns an HTTP proxy to be used for the downastream node
	Proxy() http.Handler
}
