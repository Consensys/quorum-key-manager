package keys

import (
	"context"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=keys.go -destination=mock/keys.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities2.StoreInfo, error)

	// Create a new key and stores it
	Create(ctx context.Context, id string, alg *entities2.Algorithm, attr *entities2.Attributes) (*entities2.Key, error)

	// Import an externally created key and stores it
	Import(ctx context.Context, id string, privKey []byte, alg *entities2.Algorithm, attr *entities2.Attributes) (*entities2.Key, error)

	// Get the public part of a stored key.
	Get(ctx context.Context, id string) (*entities2.Key, error)

	// List keys
	List(ctx context.Context) ([]string, error)

	// Update key tags
	Update(ctx context.Context, id string, attr *entities2.Attributes) (*entities2.Key, error)

	// Delete secret not permanently, by using Undelete() the secret can be retrieve
	Delete(ctx context.Context, id string) error

	// GetDeleted keys
	GetDeleted(ctx context.Context, id string) (*entities2.Key, error)

	// ListDeleted keys
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string) error

	// Sign from any arbitrary data using the specified key
	Sign(ctx context.Context, id string, data []byte) ([]byte, error)

	// Encrypt any arbitrary data using a specified key
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) ([]byte, error)
}
