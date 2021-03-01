package akvkeys

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
)

// Store is a key store connected to Azure Key Vault
// It delegates all crypto-operations to AKV
type Store struct {
	akv *keyvault.Client
	cfg *Config
}


func (s *AKVKeyStore) Create(ctx context.Context, id string, alg *types.Algo, attr *types.Attributes) (*types.Key, error) {
	params := keyvault.KeyCreateParameters{
		// TODO: compute keyvault.KeyCreateParameters from alg and attr
	}
	
	// Create key on AKV
	keyBundle, err := s.akv.Create(ctx, s.cfg.VaultBaseURL, id, params)
	if err != nil {
		return nil, err
	}

	key := &types.Key{
		// TODO: compute key from keyvault.KeyBundle
	}

	return key, err
}

func (s *AKVKeyStore) Sign(ctx context.Context, id string, data []byte, version int) ([]byte, error) {
	v :=  string(data)
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
