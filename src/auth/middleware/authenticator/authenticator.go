package authenticator

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=authenticator.go -destination=mock/authenticator.go -package=mock

type Authenticator interface {
	// Authenticate request

	// In case
	// - request credentials are invalid then it returns an error
	// - request credentials are valid then it returns a non UserInfo and nil error
	// - request credentials have not been passed then it returns nil, nil
	Authenticate(req *http.Request) (*types.UserInfo, error)
}

type Func func(req *http.Request) (*types.UserInfo, error)

func (f Func) Authenticate(req *http.Request) (*types.UserInfo, error) {
	return f(req)
}

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
