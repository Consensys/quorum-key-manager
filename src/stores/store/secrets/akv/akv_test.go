package akv

import (
	"context"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/akv/mocks"
	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var expectedErr = errors.AKVError("error")

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

	s.secretStore = New(s.mockVault, testutils2.NewMockLogger(ctrl))
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

	s.Run("should set a new secret successfully", func() {
		s.mockVault.EXPECT().SetSecret(gomock.Any(), id, value, attributes.Tags).Return(res, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(s.T(), attributes.Tags, secret.Tags)
		assert.Equal(s.T(), version, secret.Metadata.Version)
		assert.False(s.T(), secret.Metadata.Disabled)
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail with same error if write fails", func() {
		s.mockVault.EXPECT().SetSecret(gomock.Any(), id, value, attributes.Tags).Return(keyvault.SecretBundle{}, expectedErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(s.T(), secret)
		assert.True(s.T(), errors.IsAKVError(err))
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

	s.Run("should get a secret successfully with empty version", func() {
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(res, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(s.T(), expectedUpdatedAt, secret.Metadata.UpdatedAt)
		assert.Equal(s.T(), attributes.Tags, secret.Tags)
		assert.Equal(s.T(), version, secret.Metadata.Version)
		assert.False(s.T(), secret.Metadata.Disabled)
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail with error if bad request in response", func() {
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(keyvault.SecretBundle{}, expectedErr)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.Nil(s.T(), secret)
		assert.True(s.T(), errors.IsAKVError(err))
	})
}

func (s *akvSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	secretsList := []string{"my-secret3", "my-secret4"}

	s.Run("should list all secret ids successfully", func() {
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

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), secretsList, ids)
	})

	s.Run("should return empty list if result is nil", func() {
		s.mockVault.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return([]keyvault.SecretItem{}, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Empty(s.T(), ids)
	})

	s.Run("should fail if list fails", func() {
		s.mockVault.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return([]keyvault.SecretItem{}, expectedErr)
		ids, err := s.secretStore.List(ctx)

		assert.Nil(s.T(), ids)
		assert.True(s.T(), errors.IsAKVError(err))
	})
}

func (s *akvSecretStoreTestSuite) TestDestroy() {
	ctx := context.Background()
	id := "my-secret6"

	s.Run("should delete a secret successfully", func() {
		s.mockVault.EXPECT().PurgeDeletedSecret(gomock.Any(), id).Return(true, nil)
		err := s.secretStore.Destroy(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with NotFoundError if DeleteSecret fails with 404", func() {
		s.mockVault.EXPECT().PurgeDeletedSecret(gomock.Any(), id).Return(false, expectedErr)
		err := s.secretStore.Destroy(ctx, id)

		assert.True(s.T(), errors.IsAKVError(err))
	})
}
