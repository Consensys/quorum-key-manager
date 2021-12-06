package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
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

func (h *VaultsHandler) Register(ctx context.Context, mnfs []entities.Manifest) error {
	for _, mnf := range mnfs {
		var err error
		switch mnf.ResourceType {
		case entities.HashicorpVaultType:
			err = h.CreateHashicorp(ctx, mnf.Name, mnf.Specs)
		case entities.AzureVaultType:
			err = h.CreateAzure(ctx, mnf.Name, mnf.Specs)
		case entities.AWSVaultType:
			err = h.CreateAWS(ctx, mnf.Name, mnf.Specs)
		default:
			return errors.InvalidFormatError("invalid vault type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *VaultsHandler) CreateHashicorp(ctx context.Context, name string, specs interface{}) error {
	createReq := &types.CreateHashicorpVaultRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateHashicorp(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAzure(ctx context.Context, name string, specs interface{}) error {
	createReq := &types.CreateAzureVaultRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAzure(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (h *VaultsHandler) CreateAWS(ctx context.Context, name string, specs interface{}) error {
	createReq := &types.CreateAWSVaultRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateAWS(ctx, name, &createReq.Config, createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
