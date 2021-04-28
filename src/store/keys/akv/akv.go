package akv

import (
	"context"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
)

// Store is an implementation of key store relying on Hashicorp Vault ConsenSys secret engine
type KeyStore struct {
	client akv.SecretClient
}

var _ keys.Store = KeyStore{}

// New creates an HashiCorp key store
func New(client akv.SecretClient) *KeyStore {
	return &KeyStore{
		client:     client,
	}
}

func (k KeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Import(ctx context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Get(ctx context.Context, id, version string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) List(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	return errors.NotImplementedError
}

func (k KeyStore) Delete(ctx context.Context, id string, versions ...string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

func (k KeyStore) Destroy(ctx context.Context, id string, versions ...string) error {
	return errors.NotImplementedError
}

func (k KeyStore) Sign(ctx context.Context, id, data, version string) (string, error) {
	return "", errors.NotImplementedError
}

func (k KeyStore) Encrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.NotImplementedError
}

func (k KeyStore) Decrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.NotImplementedError
}
