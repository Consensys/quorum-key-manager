package vaults

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

type Vaults interface {
	// CreateHashicorp creates a Hashicorp Vault client
	CreateHashicorp(ctx context.Context, name string, config *entities.HashicorpConfig, allowedTenants []string, userInfo *auth.UserInfo) error

	// CreateAzure creates an AKV client
	CreateAzure(ctx context.Context, name string, config *entities.AzureConfig, allowedTenants []string, userInfo *auth.UserInfo) error

	// CreateAWS creates an AWS KMS client
	CreateAWS(ctx context.Context, name string, config *entities.AWSConfig, allowedTenants []string, userInfo *auth.UserInfo) error

	// Get gets a valut by name
	Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.Vault, error)
}
