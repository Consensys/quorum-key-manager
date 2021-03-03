package secrets

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models"
)

// Store is responsible to store secrets

// Store should be stored under file path matching regex pattern: ^[0-9a-zA-Z-]+$
type Store interface {
	// Info returns store information
	Info(context.Context) *models.StoreInfo

	// Set secret
	Set(ctx context.Context, id, value string, attr *models.Attributes) (*models.Secret, error)

	// Get a secret
	Get(ctx context.Context, id string, version string) (*models.Secret, error)

	// List secrets
	List(ctx context.Context) ([]string, error)

	// Update secret
	Update(ctx context.Context, id, newValue string, attr *models.Attributes) (*models.Secret, error)

	// Delete secret not permanently, by using Undelete the secret can be restored
	Delete(ctx context.Context, id string, versions ...int) (*models.Secret, error)

	// GetDeleted secrets
	GetDeleted(ctx context.Context, id string) (*models.Secret, error)

	// ListDeleted secrets
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanenty
	Destroy(ctx context.Context, id string, versions ...int) error
}

// Instrument allows to instrument a Store with some extra capabilities
// such as authentication, auditing, etc.
type Instrument interface {
	Apply(Store) Store
}
