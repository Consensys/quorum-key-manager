package localkeys

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/secrets"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Store is a keys.Store that uses a secrets.Store to store privateKey values
// Crypto-operations happen in the memory of the KeyStore and are not delegated to any underlying system
type Store struct {
	secrets secrets.Store
}

// New creates a new localkeys.Store
func New(secrets secrets.Store) *Store {
	return &Store{
		secrets: secrets,
	}
}

// Create a new key and stores it
func (s *Store) Create(ctx context.Context, id string, alg *types.Algo, attr *types.Attributes) (*types.Key, error) {
	switch alg.Type {
	case "ecdsa":
		// Generate key
		privKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		// Transform public key into byte
		pubKey := crypto.FromECDSAPub(privKey.Public().(*ecdsa.PublicKey))

		// Set key on the private store
		// TODO: pubkey could be stored as a metadata so we do not need to recompute it each time
		secret, err := s.secrets.Set(ctx, id, crypto.FromECDSA(privKey), attr)
		if err != nil {
			return nil, err
		}

		return &types.Key{
			PublicKey: pubKey,
			Alg:       alg,
			Attr:      secret.Attr,
			Metadata:  secret.Metadata,
		}, nil
	default:
		return nil, fmt.Errorf("not supported")
	}

}

// Sign from a digest using the specified key
func (s *Store) Sign(ctx context.Context, id string, data []byte, version int) ([]byte, error) {
	// Get secret value from secret store
	secret, err := s.secrets.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}

	// Mount secret into a private key
	privKey, err := crypto.ToECDSA(secret.Value)
	if err != nil {
		return nil, err
	}

	// Signs
	return crypto.Sign(data, privKey)
}
