package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/http/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/http/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/src/http/api/keys"
	secretsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/http/api/secrets"
	"github.com/gorilla/mux"
)

const (
	secretsPrefix  = "/secrets"
	keysPrefix     = "/keys"
	accountsPrefix = "/accounts"
	jsonRPCPrefix  = "/jsonrpc"
)

func New(bcknd core.Backend) http.Handler {
	r := mux.NewRouter()
	r.PathPrefix(secretsPrefix).Handler(stripPrefix(secretsPrefix, secretsapi.New(bcknd)))
	r.PathPrefix(keysPrefix).Handler(stripPrefix(keysPrefix, keysapi.New(bcknd)))
	r.PathPrefix(accountsPrefix).Handler(stripPrefix(accountsPrefix, accountsapi.New(bcknd)))
	r.PathPrefix(jsonRPCPrefix).Handler(stripPrefix(jsonRPCPrefix, jsonrpcapi.New(bcknd)))

	return r
}

// Modified version of http.StripPrefix() to append a tail backslash in case prefix exact match with URL.path
func stripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			if p == "" {
				p = "/"
			}

			r2.URL.Path = p
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}
