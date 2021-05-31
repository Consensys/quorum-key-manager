package manager

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/auth/types"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
)

type BaseManager struct {
	authenticators []Authenticator

	policies map[string]*Policy
}

func (a *BaseAuth) authenticate(req *http.Request) (*UserInfo, error) {
	return nil, nil
}