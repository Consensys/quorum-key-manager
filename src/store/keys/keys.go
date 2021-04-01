package keys

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

// Store is responsible to store keys and perform crypto operations

// Keys should be stored under path matching regex pattern: ^[0-9a-zA-Z-]+$
type Store interface {
	// Info returns store information
	Info(context.Context) *entities.StoreInfo

	// Create a new key and stores it
	Create(ctx context.Context, id string, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error)

	// Import an externally created key and stores it
	Import(ctx context.Context, id string, privKey []byte, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error)

	// Get the public part of a stored key.
	Get(ctx context.Context, id string, version string) (*entities.Key, error)

	// List keys
	List(ctx context.Context, count uint, skip string) (keys []*entities.Key, next string, err error)

	// Update key tags
	Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error)

	// Delete secret not parmently, by using Undelete the secret can be retrieve
	Delete(ctx context.Context, id string, versions ...string) (*entities.Key, error)

	// GetDeleted keys
	GetDeleted(ctx context.Context, id string)

	// ListDeleted keys
	ListDeleted(ctx context.Context, count uint, skip string) (keys []*entities.Key, next string, err error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string, versions ...string) error

	// Sign from a digest using the specified key
	Sign(ctx context.Context, id string, data []byte, version string) ([]byte, error)

	// Verify a signature using a specified key
	Verify(ctx context.Context, id string, data []byte) (*entities.Metadata, error)

	// Encrypt an arbitrary sequence of bytes using an encryption key that is stored in a key vault
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) (*entities.Metadata, error)
}
