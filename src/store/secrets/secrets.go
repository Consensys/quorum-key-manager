package secrets

import (
	"context"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

//go:generate mockgen -source=secrets.go -destination=mocks/secrets.go -package=mocks

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Set secret
	Set(ctx context.Context, id, value string, tags map[string]string) (*entities.Secret, error)

	// Get a secret
	Get(ctx context.Context, id string, version int) (*entities.Secret, error)

	// List secrets
	List(ctx context.Context) ([]string, error)

	// Update secret
	Refresh(ctx context.Context, id string, expirationDate time.Time) error

	// Delete secret not permanently, by using Undelete the secret can be restored
	Delete(ctx context.Context, id string, versions ...int) (*entities.Secret, error)

	// GetDeleted secrets
	GetDeleted(ctx context.Context, id string) (*entities.Secret, error)

	// ListDeleted secrets
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string, versions ...int) error
}
