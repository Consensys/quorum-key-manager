package authenticatedmanager

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/auth"
)

type Manager struct {
	mngr manager.Manager
}

func (mngr *Manager) GetSecretStore(ctx context.Context, name string) (secrets.Store, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(name); err != nil {
		return nil, err
	}

	return mngr.mngr.GetSecretStore(ctx, name)
}