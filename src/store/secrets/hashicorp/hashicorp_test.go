package hashicorp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"

	"bou.ke/monkey"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	"github.com/golang/mock/gomock"
	hashicorp "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type hashicorpSecretStoreTestSuite struct {
	suite.Suite
	mockVault   *mocks.MockVaultClient
	mountPoint  string
	secretStore secrets.Store
}

func TestHashicorpSecretStore(t *testing.T) {
	s := new(hashicorpSecretStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpSecretStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mountPoint = "secret"
	s.mockVault = mocks.NewMockVaultClient(ctrl)

	s.secretStore = New(s.mockVault, s.mountPoint)
}

func (s *hashicorpSecretStoreTestSuite) TestSet() {
	ctx := context.Background()
	id := "my-secret2"
	value := "my-value2"
	expectedPath := s.mountPoint + "/data/" + id
	attributes := testutils.FakeAttributes()
	expectedData := map[string]interface{}{
		dataLabel: map[string]interface{}{
			valueLabel: value,
			tagsLabel:  attributes.Tags,
		},
	}
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			"created_time":  "2018-03-22T02:24:06.945319214Z",
			"deletion_time": "",
			"destroyed":     false,
			"version":       json.Number("2"),
		},
	}

	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")

		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(t, attributes.Tags, secret.Tags)
		assert.Equal(t, "2", secret.Metadata.Version)
		assert.False(t, secret.Metadata.Disabled)
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusBadRequest,
		})

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
		})

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})

	// TODO: Implement specific error types and check that the function return the right error type
	s.T().Run("should fail with error if it fails to extract metadata", func(t *testing.T) {
		hashSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"created_time": "invalidTime",
				"version":      json.Number("2"),
			},
		}

		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(hashSecret, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.Error(t, err)
	})
}

func (s *hashicorpSecretStoreTestSuite) TestGet() {
	ctx := context.Background()
	id := "my-secret"
	value := "my-value"
	attributes := testutils.FakeAttributes()
	expectedPathData := s.mountPoint + "/data/" + id
	expectedPathMetadata := s.mountPoint + "/metadata/" + id

	expectedData := map[string]interface{}{
		valueLabel: value,
		tagsLabel:  attributes.Tags,
	}
	hashicorpSecretData := &hashicorp.Secret{
		Data: map[string]interface{}{
			dataLabel: expectedData,
		},
	}

	expectedMetadata := map[string]interface{}{
		"created_time":         "2018-03-22T02:24:06.945319214Z",
		"current_version":      json.Number("3"),
		"max_versions":         0,
		"oldest_version":       0,
		"updated_time":         "2018-03-22T02:36:43.986212308Z",
		"delete_version_after": "30s",
		"versions": map[string]interface{}{
			"1": map[string]interface{}{
				"created_time":  "2018-01-22T02:36:43.986212308Z",
				"deletion_time": "2018-02-22T02:36:43.986212308Z",
				"destroyed":     true,
			},
			"2": map[string]interface{}{
				"created_time":  "2018-03-22T02:36:33.954880664Z",
				"deletion_time": "",
				"destroyed":     false,
			},
			"3": map[string]interface{}{
				"created_time":  "2018-03-22T02:36:43.986212308Z",
				"deletion_time": "",
				"destroyed":     false,
			},
		},
	}
	hashicorpSecretMetadata := &hashicorp.Secret{
		Data: expectedMetadata,
	}

	s.T().Run("should get a secret successfully with empty version", func(t *testing.T) {
		expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:36:43.986212308Z")

		s.mockVault.EXPECT().Read(expectedPathData, nil).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().Read(expectedPathMetadata, nil).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(t, attributes.Tags, secret.Tags)
		assert.Equal(t, "3", secret.Metadata.Version)
		assert.False(t, secret.Metadata.Disabled)
		assert.Equal(t, secret.Metadata.CreatedAt.Add(time.Second*30), secret.Metadata.ExpireAt)
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should get a secret successfully with version", func(t *testing.T) {
		version := "2"

		s.mockVault.EXPECT().Read(expectedPathData, map[string][]string{
			versionLabel: {version},
		}).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().Read(expectedPathMetadata, nil).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(t, err)
		assert.Equal(t, secret.Metadata.Version, version)
	})

	s.T().Run("should get a secret successfully with deletion time and destroyed", func(t *testing.T) {
		version := "1"

		s.mockVault.EXPECT().Read(expectedPathData, map[string][]string{
			versionLabel: {version},
		}).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().Read(expectedPathMetadata, nil).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(t, err)
		assert.Equal(t, secret.Metadata.Version, version)
		assert.NotEmpty(t, secret.Metadata.DeletedAt)
		assert.Equal(t, secret.Metadata.DestroyedAt, secret.Metadata.DeletedAt)
	})

	s.T().Run("should fail with NotFoundError if read fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPathData, nil).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})
		s.mockVault.EXPECT().Read(expectedPathMetadata, nil).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.Nil(t, secret)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with HashicorpVaultError error if read fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPathData, nil).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})
		s.mockVault.EXPECT().Read(expectedPathMetadata, nil).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.Nil(t, secret)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/metadata"
	keys := []interface{}{"my-secret1", "my-secret2"}
	keysStr := []string{"my-secret1", "my-secret2"}

	s.T().Run("should list all secret ids successfully", func(t *testing.T) {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": keys,
			},
		}

		s.mockVault.EXPECT().List(expectedPath).Return(hashicorpSecret, nil)

		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, keysStr, ids)
	})

	s.T().Run("should return empty list if result is nil", func(t *testing.T) {
		s.mockVault.EXPECT().List(expectedPath).Return(nil, nil)

		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	s.T().Run("should fail with HashicorpVaultConnection error if write fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().List(expectedPath).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		ids, err := s.secretStore.List(ctx)

		assert.Nil(t, ids)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestRefresh() {
	ctx := context.Background()
	id := "my-secret-3"
	expectedPath := s.mountPoint + "/metadata/" + id
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{},
	}

	s.T().Run("should refresh a secret without expiration date", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{}).Return(hashicorpSecret, nil)

		err := s.secretStore.Refresh(ctx, id, "", time.Time{})

		assert.NoError(t, err)
	})

	s.T().Run("should refresh a secret with expiration date", func(t *testing.T) {
		monkey.Patch(time.Now, func() time.Time {
			return time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC)
		})
		expirationDate := time.Date(2021, 1, 1, 2, 0, 0, 0, time.UTC)

		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			deleteAfterLabel: "1h0m0s",
		}).Return(hashicorpSecret, nil)

		err := s.secretStore.Refresh(ctx, id, "", expirationDate)

		assert.NoError(t, err)
	})

	s.T().Run("should fail with NotFound error if write fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})

		err := s.secretStore.Refresh(ctx, id, "", time.Time{})

		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidFormat error if write fails with 400", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusBadRequest,
		})

		err := s.secretStore.Refresh(ctx, id, "", time.Time{})

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with InvalidParameter error if write fails with 422", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
		})

		err := s.secretStore.Refresh(ctx, id, "", time.Time{})

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with HashicorpVaultConnection error if write fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		err := s.secretStore.Refresh(ctx, id, "", time.Time{})

		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}
