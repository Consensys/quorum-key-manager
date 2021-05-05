package akv

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type akvSecretStoreTestSuite struct {
	suite.Suite
	mockVault   *mocks.MockClient
	mountPoint  string
	secretStore secrets.Store
}

func TestAkvSecretStore(t *testing.T) {
	s := new(akvSecretStoreTestSuite)
	suite.Run(t, s)
}

func (s *akvSecretStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mountPoint = "secret"
	s.mockVault = mocks.NewMockClient(ctrl)

	s.secretStore = New(s.mockVault)
}

func (s *akvSecretStoreTestSuite) TestSet() {
	ctx := context.Background()
	id := "my-secret1"
	version := "2"
	secretBundleID := id + "/" + version
	value := "my-value1"
	attributes := testutils.FakeAttributes()

	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")

	res := keyvault.SecretBundle{
		Value: &value,
		ID:    &secretBundleID,
		Attributes: &keyvault.SecretAttributes{
			Created: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedCreatedAt.UnixNano())}).x,
			Updated: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedUpdatedAt.UnixNano())}).x,
			Enabled: &(&struct{ x bool }{true}).x,
		},
		Tags: common.Tomapstrptr(attributes.Tags),
	}

	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		s.mockVault.EXPECT().SetSecret(gomock.Any(), id, value, attributes.Tags).Return(res, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(t, attributes.Tags, secret.Tags)
		assert.Equal(t, version, secret.Metadata.Version)
		assert.False(t, secret.Metadata.Disabled)
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		akvErr := autorest.DetailedError{
			Original:   fmt.Errorf("error"),
			StatusCode: 0,
		}

		s.mockVault.EXPECT().SetSecret(gomock.Any(), id, value, attributes.Tags).Return(keyvault.SecretBundle{}, akvErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.Equal(t, errors.AKVConnectionError("%v", akvErr.Original), err)
	})
}

func (s *akvSecretStoreTestSuite) TestGet() {
	ctx := context.Background()
	id := "my-secret2"
	version := "2"
	secretBundleID := id + "/" + version
	value := "my-value2"
	attributes := testutils.FakeAttributes()

	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2018-03-23T02:24:06.945319214Z")

	res := keyvault.SecretBundle{
		Value: &value,
		ID:    &secretBundleID,
		Attributes: &keyvault.SecretAttributes{
			Created: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedCreatedAt.UnixNano())}).x,
			Updated: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedUpdatedAt.UnixNano())}).x,
			Enabled: &(&struct{ x bool }{true}).x,
		},
		Tags: common.Tomapstrptr(attributes.Tags),
	}

	s.T().Run("should get a secret successfully with empty version", func(t *testing.T) {
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(res, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(t, expectedUpdatedAt, secret.Metadata.UpdatedAt)
		assert.Equal(t, attributes.Tags, secret.Tags)
		assert.Equal(t, version, secret.Metadata.Version)
		assert.False(t, secret.Metadata.Disabled)
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with error if bad request in response", func(t *testing.T) {
		akvErr := autorest.DetailedError{
			Original:   fmt.Errorf("error"),
			StatusCode: http.StatusBadRequest,
		}

		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(keyvault.SecretBundle{}, akvErr)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.Nil(t, secret)
		assert.Equal(t, errors.InvalidFormatError("%v", akvErr.Original), err)
	})
}

func (s *akvSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	secretsList := []string{"my-secret3", "my-secret4"}

	s.T().Run("should list all secret ids successfully", func(t *testing.T) {
		items := []keyvault.SecretItem{
			{
				ID: &(&struct{ x string }{"https://test.dns/secrets/my-secret3"}).x,
			},
			{
				ID: &(&struct{ x string }{"https://test.dns/secrets/my-secret4"}).x,
			},
		}
		result := keyvault.SecretListResult{
			Value: &items,
		}
		list := keyvault.NewSecretListResultPage(result, nil).Values()

		s.mockVault.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return(list, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, secretsList, ids)
	})

	s.T().Run("should return empty list if result is nil", func(t *testing.T) {
		s.mockVault.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return([]keyvault.SecretItem{}, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Nil(t, ids)
	})

	s.T().Run("should fail if list fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		akvErr := autorest.DetailedError{
			Original:   expectedErr,
			StatusCode: http.StatusNotFound,
		}

		s.mockVault.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return([]keyvault.SecretItem{}, akvErr)
		ids, err := s.secretStore.List(ctx)

		assert.Nil(t, ids)
		assert.True(t, errors.IsNotFoundError(err))
		assert.Equal(t, errors.NotFoundError("%v", expectedErr), err)
	})
}

func (s *akvSecretStoreTestSuite) TestRefresh() {
	ctx := context.Background()
	id := "my-secret5"
	version := "2"

	expectedExpirationDate, _ := time.Parse(time.RFC3339, "2021-03-22T02:24:06.945319214Z")
	s.T().Run("should refresh a secret with expiration date", func(t *testing.T) {
		s.mockVault.EXPECT().UpdateSecret(gomock.Any(), id, version, expectedExpirationDate).Return(keyvault.SecretBundle{}, nil)
		err := s.secretStore.Refresh(ctx, id, version, expectedExpirationDate)
		assert.NoError(t, err)
	})

	s.T().Run("should fail if UpdateSecret fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		akvErr := autorest.DetailedError{
			Original:   expectedErr,
			StatusCode: http.StatusNotFound,
		}

		s.mockVault.EXPECT().UpdateSecret(gomock.Any(), id, version, expectedExpirationDate).Return(keyvault.SecretBundle{}, akvErr)
		err := s.secretStore.Refresh(ctx, id, version, expectedExpirationDate)

		assert.True(t, errors.IsNotFoundError(err))
		assert.Equal(t, errors.NotFoundError("%v", expectedErr), err)
	})
}

func (s *akvSecretStoreTestSuite) TestDestroy() {
	ctx := context.Background()
	id := "my-secret6"

	s.T().Run("should delete a secret successfully", func(t *testing.T) {
		s.mockVault.EXPECT().PurgeDeletedSecret(gomock.Any(), id).Return(true, nil)
		err := s.secretStore.Destroy(ctx, id)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with NotFoundError if DeleteSecret fails with 404", func(t *testing.T) {
		akvErr := errors.NotFoundError("not found")

		s.mockVault.EXPECT().PurgeDeletedSecret(gomock.Any(), id).Return(false, akvErr)
		err := s.secretStore.Destroy(ctx, id)

		assert.True(t, errors.IsNotFoundError(err))
	})
}
