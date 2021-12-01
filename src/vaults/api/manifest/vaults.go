package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/vaults"
)

type VaultsHandler struct {
	vaults vaults.Vaults
}

func NewVaultsHandler(vaultsService vaults.Vaults) *VaultsHandler {
	return &VaultsHandler{
		vaults: vaultsService,
	}
}

func (h *VaultsHandler) CreateHashicorp(ctx context.Context, config interface{}) error {
	createReq := &types.CreateHashicorpVaultRequest{}
	err := json.UnmarshalJSON(config, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateHashicorp(ctx, createReq.Name, &createReq.Config)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAzure(ctx context.Context, params interface{}) error {
	createReq := &types.CreateAzureVaultRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAzure(ctx, createReq.Name, &createReq.Config)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAWS(ctx context.Context, params interface{}) error {
	createReq := &types.CreateAWSVaultRequest{}
	err := json.UnmarshalJSON(params, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAWS(ctx, createReq.Name, &createReq.Config)
	if err != nil {
		return err
	}

	return nil
}
