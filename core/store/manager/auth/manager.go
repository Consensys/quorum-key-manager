package authenticatedmanager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/auth"
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/secrets"
)

type Manager struct {
	mngr manager.Manager
}

func Wrap(mngr manager.Manager) *Manager {
	return &Manager{
		mngr: mngr,
	}
}

func (mngr *Manager) Load(ctx context.Context, mnfsts ...*manifestloader.Message) {
	// TODO: at this stage, store creation is only handled by system so no need for permission checks
	mngr.mngr.Load(ctx, mnfsts...)
}

func (mngr *Manager) GetSecretStore(ctx context.Context, name string) (secrets.Store, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(name); err != nil {
		return nil, err
	}

	return mngr.mngr.GetSecretStore(ctx, name)
}

func (mngr *Manager) GetKeyStore(ctx context.Context, name string) (keys.Store, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(name); err != nil {
		return nil, err
	}

	return mngr.mngr.GetKeyStore(ctx, name)
}

func (mngr *Manager) GetAccountStore(ctx context.Context, name string) (accounts.Store, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(name); err != nil {
		return nil, err
	}

	return mngr.mngr.GetAccountStore(ctx, name)
}
