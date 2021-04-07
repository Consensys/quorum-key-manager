package akv

import (
	"context"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

type SecretStore struct {
	client akv.Client
}

func New(client akv.Client) *SecretStore {
	return &SecretStore{
		client: client,
	}
}

func (s *SecretStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

func (s *SecretStore) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	params := keyvault.SecretSetParameters{
		Value: &value,
		Tags:  tomapstrptr(attr.Tags),
	}

	res, err := s.client.SetSecret(ctx, id, params)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return parseSecretBundle(res), nil
}

func (s *SecretStore) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	res, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return parseSecretBundle(res), nil
}

func (s *SecretStore) List(ctx context.Context) ([]string, error) {
	res, err := s.client.GetSecrets(ctx, nil)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	if len(res.Values()) == 0 {
		return nil, nil
	}
	var list []string
	for _, secret := range res.Values() {
		// path.Base to only retrieve the secretName instead of https://<vaultName>.vault.azure.net/secrets/<secretName>
		// See listSecrets function in https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/keyvault/examples/go-keyvault-msi-example.go
		list = append(list, path.Base(*secret.ID))
	}
	return list, nil
}

func (s *SecretStore) Refresh(ctx context.Context, id, version string, expirationDate time.Time) error {
	expires := date.NewUnixTimeFromNanoseconds(expirationDate.UnixNano())
	params := keyvault.SecretUpdateParameters{
		SecretAttributes: &keyvault.SecretAttributes{
			Expires: &expires,
		},
	}

	_, err := s.client.UpdateSecret(ctx, id, version, params)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (s *SecretStore) Delete(ctx context.Context, id string, versions ...string) (*entities.Secret, error) {

	return nil, errors.NotImplementedError
}

func (s *SecretStore) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	return nil, errors.NotImplementedError
}

func (s *SecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

func (s *SecretStore) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

func (s *SecretStore) Destroy(ctx context.Context, id string, versions ...string) error {
	return errors.NotImplementedError
}
