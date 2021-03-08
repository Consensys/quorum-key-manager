package akvkeys

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

// Store is a key store connected to Azure Key Vault
// It delegates all crypto-operations to AKV
type Store struct {
	akv *keyvault.Client
	cfg *Config
}

// New create a new Key Store connected to Azure Key Vault
// It delegates all crypto-operations to AKV
func New(cfg *Config) (*Store, error) {
	akv := keyvault.New()

	// TODO: prepare client from cfg

	return &Store{
		cfg: cfg,
		akv: akv,
	}, nil
}

// Create a new key and stores it
func (s *AKVKeyStore) Create(ctx context.Context, id string, alg *models.Algo, attr *models.Attributes) (*models.Key, error) {
	params := keyvault.KeyCreateParameters{
		// TODO: compute keyvault.KeyCreateParameters from alg and attr
	}

	// Create key on AKV
	keyBundle, err := s.akv.Create(ctx, s.cfg.VaultBaseURL, id, params)
	if err != nil {
		return nil, err
	}

	key := &models.Key{
		// TODO: compute key from keyvault.KeyBundle
	}

	return key, err
}

// Sign from a digest using the specified key
func (s *AKVKeyStore) Sign(ctx context.Context, id string, data []byte, version int) ([]byte, error) {
	v := string(data)
	params := keyvault.KeySignParameters{
		// TODO: compute keyvault.KeySignParameters
		Value: &v,
	}

	keyOpRes, err := s.akv.Sign(ctx, s.cfg.VaultBaseURL, id, version, params)
	if err != nil {
		return nil, err
	}

	return byte(*keyOpRes.Result), err
}
