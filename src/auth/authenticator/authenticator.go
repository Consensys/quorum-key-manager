package authenticator

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

type Authenticator interface {
	// Authenticate request

	// In case
	// - request credentials are invalid then it returns an error
	// - request credentials are valid then it returns a non UserInfo and nil error
	// - request credentials have not been passed then it returns nil, nil
	Authenticate(req *http.Request) (*types.UserInfo, error)
}
