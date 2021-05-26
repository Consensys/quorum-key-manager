package auth

import (
	"context"
	"net/http"
)

type Auth interface {
	Load(context.Context)

	// Authenticate request
	Authenticate(req *http.Request) error
}

type BaseAuth struct {
	authenticators []Authenticator

	policies map[string]*Policy
}

func (a *BaseAuth) authenticate(req *http.Request) (*UserInfo, error) {
	return nil, nil
}

func (a *BaseAuth) impersonate(req *http.Request, reqCtx *RequestContext) error {
	return nil
}
