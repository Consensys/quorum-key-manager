package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type StoresHandler struct {
	stores   stores.Stores
	userInfo *auth.UserInfo
}

func NewStoresHandler(storesConnector stores.Stores) *StoresHandler {
	return &StoresHandler{
		stores:   storesConnector,
		userInfo: auth.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *StoresHandler) Register(ctx context.Context, specs interface{}) error {
	createReq := &types.CreateStoreRequest{}
	err := json.UnmarshalJSON(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	switch createReq.StoreType {
	case entities.SecretStoreType:
		return h.CreateSecret(ctx, createReq.Params)
	case entities.KeyStoreType:
		return h.CreateSecret(ctx, createReq.Params)
	case entities.EthereumStoreType:
		return h.CreateSecret(ctx, createReq.Params)
	default:
		return errors.InvalidFormatError("invalid store type")
	}
}

func (h *StoresHandler) CreateSecret(ctx context.Context, params interface{}) error {
	createReq := &types.CreateSecretStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateSecret(ctx, createReq.Name, createReq.Vault, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateKey(ctx context.Context, params interface{}) error {
	createReq := &types.CreateKeyStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateKey(ctx, createReq.Name, createReq.Vault, createReq.SecretStore, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *StoresHandler) CreateEthereum(ctx context.Context, params interface{}) error {
	createReq := &types.CreateEthereumStoreRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.stores.CreateEthereum(ctx, createReq.Name, createReq.KeyStore, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
