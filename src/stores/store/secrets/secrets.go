package secrets

import (
	"context"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=secrets.go -destination=mock/secrets.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities2.StoreInfo, error)

	// Set secret
	Set(ctx context.Context, id, value string, attr *entities2.Attributes) (*entities2.Secret, error)

	// Get a secret
	Get(ctx context.Context, id string, version string) (*entities2.Secret, error)

	// List secrets
	List(ctx context.Context) ([]string, error)

	// Delete secret not permanently, by using Undelete the secret can be restored
	Delete(ctx context.Context, id string) error

	// GetDeleted secrets
	GetDeleted(ctx context.Context, id string) (*entities2.Secret, error)

	// ListDeleted secrets
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string) error
}
