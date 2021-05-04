package akv

import (
	"context"
	"path"
	"time"

	akvclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

type Store struct {
	client akv.SecretClient
}

var _ secrets.Store = &Store{}

func New(client akv.SecretClient) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	res, err := s.client.SetSecret(ctx, id, value, attr.Tags)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseSecretBundle(res), nil
}

func (s *Store) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	res, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseSecretBundle(res), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	items, err := s.client.GetSecrets(ctx, 0)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	var list []string
	for _, secret := range items {
		// path.Base to only retrieve the secretName instead of https://<vaultName>.vault.azure.net/secrets/<secretName>
		// See listSecrets function in https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/keyvault/examples/go-keyvault-msi-example.go
		list = append(list, path.Base(*secret.ID))
	}
	return list, nil
}

func (s *Store) Refresh(ctx context.Context, id, version string, expirationDate time.Time) error {
	_, err := s.client.UpdateSecret(ctx, id, version, expirationDate)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, id string) (*entities.Secret, error) {

	return nil, errors.ErrNotImplemented
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	_, err := s.client.DeleteSecret(ctx, id)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}
