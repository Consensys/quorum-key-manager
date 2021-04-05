//nolint
package akvkeys

import (
	"context"

	// "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
)

// Store is a key store connected to Azure Key Vault
// It delegates all crypto-operations to AKV
type store struct {
	// akv *keyvault.BaseClient
	cfg *Config
}

// New create a new Key Store connected to Azure Key Vault
// It delegates all crypto-operations to AKV
func New(cfg *Config) (keys.Store, error) {
	// akv := keyvault.New()

	// TODO: prepare client from cfg

	return &store{
		cfg: cfg,
		// akv: &akv,
	}, nil
}

// Create a new key and stores it
func (s *store) Create(ctx context.Context, id string, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error) {
	panic("implement me")
	// params := keyvault.KeyCreateParameters{
	// 	// TODO: compute keyvault.KeyCreateParameters from alg and attr
	// }
	//
	// // Create key on AKV
	// _, err := s.akv.CreateKey(ctx, s.cfg.VaultBaseURL, id, params)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// key := &entities.Key{
	// 	// TODO: compute key from keyvault.KeyBundle
	// }
	//
	// return key, err
}

// Sign from a digest using the specified key
func (s *store) Sign(ctx context.Context, id string, data []byte, version string) ([]byte, error) {
	panic("implement me")
	// v := string(data)
	// params := keyvault.KeySignParameters{
	// 	// TODO: compute keyvault.KeySignParameters
	// 	Value: &v,
	// }
	//
	// keyOpRes, err := s.akv.Sign(ctx, s.cfg.VaultBaseURL, id, string(version), params)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return []byte(*keyOpRes.Result), err
}

func (s *store) Info(context.Context) *entities.StoreInfo {
	panic("implement me")
}

func (s *store) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error) {
	panic("implement me")
}

func (s *store) Get(ctx context.Context, id string, version string) (*entities.Key, error) {
	panic("implement me")
}

func (s *store) List(ctx context.Context, count uint, skip string) (keys []*entities.Key, next string, err error) {
	panic("implement me")
}

func (s *store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	panic("implement me")
}

func (s *store) Delete(ctx context.Context, id string, versions ...string) (*entities.Key, error) {
	panic("implement me")
}

func (s *store) GetDeleted(ctx context.Context, id string) {
	panic("implement me")
}

func (s *store) ListDeleted(ctx context.Context, count uint, skip string) (keys []*entities.Key, next string, err error) {
	panic("implement me")
}

func (s *store) Undelete(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *store) Destroy(ctx context.Context, id string, versions ...string) error {
	panic("implement me")
}

func (s *store) Verify(ctx context.Context, id string, data []byte) (*entities.Metadata, error) {
	panic("implement me")
}

func (s *store) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	panic("implement me")
}

func (s *store) Decrypt(ctx context.Context, id string, data []byte) (*entities.Metadata, error) {
	panic("implement me")
}
