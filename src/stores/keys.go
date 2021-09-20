package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

//go:generate mockgen -source=keys.go -destination=mock/keys.go -package=mock

type KeyStore interface {
	// Create creates a new key and stores it
	Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error)

	// Import imports an externally created key and stores it
	Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error)

	// Get gets the public part of a stored key.
	Get(ctx context.Context, id string) (*entities.Key, error)

	// List lists keys
	List(ctx context.Context, limit, offset uint64) ([]string, error)

	// Update updates key tags
	Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error)

	// Delete soft-deletes a key
	Delete(ctx context.Context, id string) error

	// GetDeleted gets a deleted key
	GetDeleted(ctx context.Context, id string) (*entities.Key, error)

	// ListDeleted lists deleted keys
	ListDeleted(ctx context.Context, limit, offset uint64) ([]string, error)

	// Restore restores a previously deleted secret
	Restore(ctx context.Context, id string) error

	// Destroy destroys a key permanently
	Destroy(ctx context.Context, id string) error

	// Sign from any arbitrary data using the specified key
	Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error)

	// Encrypt encrypts any arbitrary data using a specified key
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt decrypts a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) ([]byte, error)
}
