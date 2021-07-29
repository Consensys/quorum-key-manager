package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=secrets.go -destination=mock/secrets.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Set secret
	Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error)

	// Get a secret
	Get(ctx context.Context, id string, version string) (*entities.Secret, error)

	// List secrets
	List(ctx context.Context) ([]string, error)

	// Delete secret not permanently, it can be restored
	Delete(ctx context.Context, id string, version string) error

	// GetDeleted secrets
	GetDeleted(ctx context.Context, id, version string) (*entities.Secret, error)

	// ListDeleted secrets
	ListDeleted(ctx context.Context) ([]string, error)

	// Restore a previously deleted secret
	Restore(ctx context.Context, id, version string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id, version string) error
}
