package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type VaultsHandler struct {
	vaults stores.Vaults
}

func NewVaultsHandler(vaultsConnector stores.Vaults) *VaultsHandler {
	return &VaultsHandler{
		vaults: vaultsConnector,
	}
}

func (h *VaultsHandler) Register(ctx context.Context, specs interface{}) error {
	createReq := &types.CreateVaultRequest{}
	err := json.UnmarshalJSON(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	switch createReq.VaultType {
	case entities.HashicorpVaultType:
		return h.CreateHashicorp(ctx, createReq.Params)
	case entities.AzureVaultType:
		return h.CreateAzure(ctx, createReq.Params)
	case entities.AWSVaultType:
		return h.CreateAWS(ctx, createReq.Params)
	default:
		return errors.InvalidFormatError("invalid vault type")
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
