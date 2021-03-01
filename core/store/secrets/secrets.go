package secrets

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

// Store is responsible to store secrets

// Store should be stored under file path matching regex pattern: ^[0-9a-zA-Z-]+$
type Store interface {
	// Info returns store information
	Info(context.Context) *types.StoreInfo

	// Set secret
	Set(ctx context.Context, id string, value []byte, attr *types.Attributes) (*types.Metadata, error)

	// Get a secret
	Get(ctx context.Context, id string, version int) (*types.Secret, error)

	// List secrets
	List(ctx context.Context, count uint, skip string) (secrets []*types.Secret, next string, err error)

	// Update secret
	Update(ctx context.Context, id string, attr *types.Attributes)

	// Delete secret not permanently, by using Undelete the secret can be restored
	Delete(ctx context.Context, id string, versions ...int) (*types.Secret, error)

	// GetDeleted secrets
	GetDeleted(ctx context.Context, id string) (*types.Secret, error)

	// ListDeleted secrets
	ListDeleted(ctx context.Context, count uint, skip string) (secrets []*types.Secret, next string, err error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanenty
	Destroy(ctx context.Context, id string, versions ...int) error
}
