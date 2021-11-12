package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

//go:generate mockgen -source=vaults.go -destination=mock/vaults.go -package=mock

type Vaults interface {
	// CreateHashicorp creates a Hashicorp Vault client
	CreateHashicorp(ctx context.Context, name string, config *entities.HashicorpConfig) error

	// CreateAzure creates an AKV client
	CreateAzure(ctx context.Context, name string, config *entities.AzureConfig) error

	// CreateAWS creates an AWS KMS client
	CreateAWS(ctx context.Context, name string, config *entities.AWSConfig) error

	// Get gets a valut by name
	Get(ctx context.Context, name string) (*entities.Vault, error)
}
