package akv

import (
	"context"
	"path"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type Store struct {
	client akv.SecretClient
	logger log.Logger
}

var _ stores.SecretStore = &Store{}

func New(client akv.SecretClient, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	res, err := s.client.SetSecret(ctx, id, value, attr.Tags)
	if err != nil {
		errMessage := "failed to create AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseSecretBundle(&res), nil
}

func (s *Store) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	res, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		errMessage := "failed to get AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseSecretBundle(&res), nil
}

func (s *Store) List(ctx context.Context, _, _ int) ([]string, error) {
	items, err := s.client.ListSecrets(ctx, 0)
	if err != nil {
		errMessage := "failed to list AKV secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var list = []string{}
	for _, secret := range items {
		// path.Base to only retrieve the secretName instead of https://<vaultName>.vault.azure.net/secrets/<secretName>
		// See listSecrets function in https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/keyvault/examples/go-keyvault-msi-example.go
		list = append(list, path.Base(*secret.ID))
	}

	return list, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.client.DeleteSecret(ctx, id)
	if err != nil {
		errMessage := "failed to delete AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	res, err := s.client.GetDeletedSecret(ctx, id)
	if err != nil {
		errMessage := "failed to get deleted AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseDeletedSecretBundle(&res), nil
}

func (s *Store) ListDeleted(ctx context.Context, _, _ int) ([]string, error) {
	items, err := s.client.ListDeletedSecrets(ctx, 0)
	if err != nil {
		errMessage := "failed to list deleted AKV secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var list = []string{}
	for _, secret := range items {
		// path.Base to only retrieve the secretName instead of https://<vaultName>.vault.azure.net/secrets/<secretName>
		// See listSecrets function in https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/keyvault/examples/go-keyvault-msi-example.go
		list = append(list, path.Base(*secret.ID))
	}

	return list, nil
}

func (s *Store) Restore(ctx context.Context, id string) error {
	_, err := s.client.RecoverSecret(ctx, id)
	if err != nil {
		errMessage := "failed to restore AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	_, err := s.client.PurgeDeletedSecret(ctx, id)
	if err != nil {
		errMessage := "failed to permanently delete AKV secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
