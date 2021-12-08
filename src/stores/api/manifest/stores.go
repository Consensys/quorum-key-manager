package manifest

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type StoresHandler struct {
	stores   stores.Stores
	userInfo *authtypes.UserInfo
}

func NewStoresHandler(storesService stores.Stores) *StoresHandler {
	return &StoresHandler{
		stores:   storesService,
		userInfo: authtypes.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *StoresHandler) Register(ctx context.Context, mnfs []entities2.Manifest) error {
	for _, mnf := range mnfs {
		var err error
		switch mnf.ResourceType {
		case entities.SecretStoreType:
			err = h.CreateSecret(ctx, mnf.Name, mnf.AllowedTenants, mnf.Specs)
		case entities.KeyStoreType:
			err = h.CreateKey(ctx, mnf.Name, mnf.AllowedTenants, mnf.Specs)
		case entities.EthereumStoreType:
			err = h.CreateEthereum(ctx, mnf.Name, mnf.AllowedTenants, mnf.Specs)
		default:
			err = errors.InvalidFormatError("invalid store type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *StoresHandler) CreateSecret(ctx context.Context, name string, allowedTenants []string, specs interface{}) error {
	createReq := &types.CreateSecretStoreRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateSecret(ctx, name, createReq.Vault, allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateKey(ctx context.Context, name string, allowedTenants []string, specs interface{}) error {
	createReq := &types.CreateKeyStoreRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateKey(ctx, name, createReq.Vault, createReq.SecretStore, allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateEthereum(ctx context.Context, name string, allowedTenants []string, specs interface{}) error {
	createReq := &types.CreateEthereumStoreRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateEthereum(ctx, name, createReq.KeyStore, allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
