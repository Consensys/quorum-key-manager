package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
)

type StoresHandler struct {
	stores   stores.Stores
	userInfo *entities.UserInfo
}

func NewStoresHandler(storesConnector stores.Stores) *StoresHandler {
	return &StoresHandler{
		stores:   storesConnector,
		userInfo: entities.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *StoresHandler) CreateSecret(ctx context.Context, name string, params interface{}) error {
	createReq := &types.CreateSecretStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateSecret(ctx, name, createReq.Vault, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateKey(ctx context.Context, name string, params interface{}) error {
	createReq := &types.CreateKeyStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateKey(ctx, name, createReq.Vault, createReq.SecretStore, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateEthereum(ctx context.Context, name string, params interface{}) error {
	createReq := &types.CreateEthereumStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateEthereum(ctx, name, createReq.KeyStore, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
