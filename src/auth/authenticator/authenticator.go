package authenticator

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=authenticator.go -destination=mock/authenticator.go -package=mock

type Authenticator interface {
	// Authenticate MUST
	// - return an error, if request credentials are present but invalid
	// - returns a non nil UserInfo and nil error, if request credentials are valid
	// - returns nil, nil, if request credentials are missing
	Authenticate(req *http.Request) (*types.UserInfo, error)
}

type Func func(req *http.Request) (*types.UserInfo, error)

func (f Func) Authenticate(req *http.Request) (*types.UserInfo, error) {
	return f(req)
}

// First combines authenticators

// First executes authenticators in sequence until one authenticator accepts
// or rejects the request
func First(authenticators ...Authenticator) Authenticator {
	return Func(func(req *http.Request) (*types.UserInfo, error) {
		for _, auth := range authenticators {
			info, err := auth.Authenticate(req)
			if info != nil || err != nil {
				return info, err
			}
		}
		return nil, nil
	})
}
