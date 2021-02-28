package keys

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/store/types"
)

// Store is responsible to store keys and perform crypto operations

// Keys should be stored under path matching regex pattern: ^[0-9a-zA-Z-]+$
type Store interface {
	// Info returns store information
	Info(context.Context) *types.StoreInfo

	// Create a new key and stores it
	Create(ctx context.Context, id string, alg *types.Algo, attr *types.Attributes) (*types.Key, error)

	// Import an externally created key and stores it
	Import(ctx context.Context, id string, privKey []byte, alg *types.Algo, attr *types.Attributes) (*types.Key, error)

	// Get the public part of a stored key.
	Get(ctx context.Context, id string, version int) (*types.Key, error)

	// List keys
	List(ctx context.Context, count uint, skip string) (keys []*types.Key, next string, err error)

	// Update key tags
	Update(ctx context.Context, id string, attr *types.Attributes) (*types.Key, error)

	// Delete secret not parmently, by using Undelete the secret can be retrieve
	Delete(ctx context.Context, id string, versions ...int) (*types.Key, error)

	// GetDeleted keys
	GetDeleted(ctx context.Context, id string)

	// ListDeleted keys
	ListDeleted(ctx context.Context, count uint, skip string) (keys []*types.Key, next string, err error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanenty
	Destroy(ctx context.Context, id string, versions ...int) error

	// Sign from a digest using the specified key
	Sign(ctx context.Context, id string, data []byte) ([]byte, error)

	// Verify a signature using a specified key
	Verify(ctx context.Context, id string, data []byte) (*types.Metadata, error)

	// Encrypt an arbitrary sequence of bytes using an encryption key that is stored in a key vault
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) (*types.Metadata, error)
}
