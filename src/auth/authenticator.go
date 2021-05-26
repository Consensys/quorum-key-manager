package auth

import "net/http"

type Authenticator interface {
	Authenticate(req *http.Request) (*UserInfo, error)
}

// TODO: implement AuthenticatorTLS
type AuthenticatorTLS struct{}

// TODO: implement AuthenticatorAPIKey
type AuthenticatorAPIKey struct{}

// TODO: implement AuthenticatorOIDC
type AuthenticatorOIDC struct{}
