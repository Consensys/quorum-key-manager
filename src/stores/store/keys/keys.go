package keys

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=keys.go -destination=mock/keys.go -package=mock

type Store interface {
	// Create a new key and stores it
	Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*models.Key, error)

	// Import imports an existing key and stores it
	Import(ctx context.Context, id string, importedPrivKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*models.Key, error)

	// Update updates the tags of a key
	Update(ctx context.Context, id string, attr *entities.Attributes) (*models.Key, error)

	// Delete secret not permanently, by using Undelete() the secret can be retrieved
	Delete(ctx context.Context, id string) error

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string) error

	// Sign from any arbitrary data using the specified key
	Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error)

	// Encrypt any arbitrary data using a specified key
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) ([]byte, error)
}
