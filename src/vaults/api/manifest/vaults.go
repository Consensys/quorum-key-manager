package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/vaults"
)

type VaultsHandler struct {
	vaults   vaults.Vaults
	userInfo *auth.UserInfo
}

func NewVaultsHandler(vaultsService vaults.Vaults) *VaultsHandler {
	return &VaultsHandler{
		vaults:   vaultsService,
		userInfo: auth.NewWildcardUser(),
	}
}

func (h *VaultsHandler) CreateHashicorp(ctx context.Context, name string, config interface{}) error {
	createReq := &types.CreateHashicorpVaultRequest{}
	err := json.UnmarshalJSON(config, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateHashicorp(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAzure(ctx context.Context, name string, config interface{}) error {
	createReq := &types.CreateAzureVaultRequest{}
	err := json.UnmarshalJSON(config, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAzure(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAWS(ctx context.Context, name string, config interface{}) error {
	createReq := &types.CreateAWSVaultRequest{}
	err := json.UnmarshalJSON(config, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAWS(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
