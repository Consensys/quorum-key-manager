package hashicorp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/mocks"
	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/stores"
	dbmocks "github.com/consensys/quorum-key-manager/src/stores/database/mock"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/golang/mock/gomock"
	hashicorp "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var expectedErr = errors.HashicorpVaultError("error")

type hashicorpSecretStoreTestSuite struct {
	suite.Suite
	mockVault   *mocks.MockKvv2Client
	secretStore stores.SecretStore
	mockDB      *dbmocks.MockSecrets
}

func TestHashicorpSecretStore(t *testing.T) {
	s := new(hashicorpSecretStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpSecretStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockVault = mocks.NewMockKvv2Client(ctrl)
	s.mockDB = dbmocks.NewMockSecrets(ctrl)

	s.secretStore = New(s.mockVault, s.mockDB, testutils2.NewMockLogger(ctrl))
}

func (s *hashicorpSecretStoreTestSuite) TestSet() {
	ctx := context.Background()
	id := "my-secret2"
	value := "my-value2"
	version := "3"
	attributes := testutils.FakeAttributes()
	expectedWriteData := map[string]interface{}{
		valueLabel: value,
		tagsLabel:  attributes.Tags,
	}

	expectedData := map[string]interface{}{
		valueLabel: value,
		tagsLabel: map[string]interface{}{
			"tag1": attributes.Tags["tag1"],
			"tag2": attributes.Tags["tag2"],
		},
	}
	hashicorpSecretData := &hashicorp.Secret{
		Data: map[string]interface{}{
			"data": expectedData,
		},
	}
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			"created_time":  "2018-03-22T02:24:06.945319214Z",
			"deletion_time": "",
			"destroyed":     false,
			"version":       json.Number(version),
		},
	}
	expectedMetadata := map[string]interface{}{
		"created_time":         "2018-03-22T02:24:06.945319214Z",
		"current_version":      json.Number(version),
		"max_versions":         0,
		"oldest_version":       0,
		"updated_time":         "2018-03-22T02:36:43.986212308Z",
		"delete_version_after": "0s",
		"versions": map[string]interface{}{
			"1": map[string]interface{}{
				"created_time":  "2018-01-22T02:36:43.986212308Z",
				"deletion_time": "",
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

	s.Run("should set a new secret successfully", func() {
		expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:36:43.986212308Z")

		s.mockVault.EXPECT().SetSecret(id, expectedWriteData).Return(hashicorpSecret, nil)
		s.mockVault.EXPECT().ReadData(id, gomock.Any()).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().ReadMetadata(id).Return(hashicorpSecretMetadata, nil)

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
		s.mockVault.EXPECT().SetSecret(id, expectedWriteData).Return(nil, expectedErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(s.T(), secret)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestGet() {
	ctx := context.Background()
	id := "my-get-secret"
	value := "my-value"
	attributes := testutils.FakeAttributes()

	expectedData := map[string]interface{}{
		valueLabel: value,
		tagsLabel: map[string]interface{}{
			"tag1": attributes.Tags["tag1"],
			"tag2": attributes.Tags["tag2"],
		},
	}
	hashicorpSecretData := &hashicorp.Secret{
		Data: map[string]interface{}{
			"data": expectedData,
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

	s.Run("should get a secret successfully with empty version", func() {
		expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:36:43.986212308Z")

		s.mockVault.EXPECT().ReadData(id, nil).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().ReadMetadata(id).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), expectedCreatedAt, secret.Metadata.CreatedAt)
		assert.Equal(s.T(), attributes.Tags, secret.Tags)
		assert.Equal(s.T(), "3", secret.Metadata.Version)
		assert.False(s.T(), secret.Metadata.Disabled)
		assert.Equal(s.T(), secret.Metadata.CreatedAt.Add(time.Second*30), secret.Metadata.ExpireAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
	})

	s.Run("should get a secret successfully with version", func() {
		version := "2"

		s.mockVault.EXPECT().ReadData(id, map[string][]string{versionLabel: {version}}).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().ReadMetadata(id).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), secret.Metadata.Version, version)
	})

	s.Run("should get a secret successfully with deletion time and destroyed", func() {
		version := "1"

		s.mockVault.EXPECT().ReadData(id, map[string][]string{
			versionLabel: {version},
		}).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().ReadMetadata(id).Return(hashicorpSecretMetadata, nil)

		secret, err := s.secretStore.Get(ctx, id, version)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), secret.Metadata.Version, version)
		assert.NotEmpty(s.T(), secret.Metadata.DeletedAt)
	})

	s.Run("should fail with same error if read data fails", func() {
		s.mockVault.EXPECT().ReadData(id, nil).Return(nil, expectedErr)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.Nil(s.T(), secret)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})

	s.Run("should fail with same error if read metadata fails", func() {
		s.mockVault.EXPECT().ReadData(id, nil).Return(hashicorpSecretData, nil)
		s.mockVault.EXPECT().ReadMetadata(id).Return(nil, expectedErr)

		secret, err := s.secretStore.Get(ctx, id, "")

		assert.Nil(s.T(), secret)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	keys := []interface{}{"my-secret1", "my-secret2"}
	keysStr := []string{"my-secret1", "my-secret2"}

	s.Run("should list all secret ids successfully", func() {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": keys,
			},
		}

		s.mockVault.EXPECT().ListSecrets().Return(hashicorpSecret, nil)

		ids, err := s.secretStore.List(ctx, 0, 0)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), keysStr, ids)
	})

	s.Run("should return empty list if result is nil", func() {
		s.mockVault.EXPECT().ListSecrets().Return(nil, nil)

		ids, err := s.secretStore.List(ctx, 0, 0)

		assert.NoError(s.T(), err)
		assert.Empty(s.T(), ids)
	})

	s.Run("should fail with same error if ListSecrets fails", func() {
		s.mockVault.EXPECT().ListSecrets().Return(nil, expectedErr)

		ids, err := s.secretStore.List(ctx, 0, 0)

		assert.Empty(s.T(), ids)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestDelete() {
	ctx := context.Background()
	id := "my-deleted-secret"
	versions := []string{"1", "2", "3"}

	s.Run("should delete secret by id successfully", func() {
		data := map[string][]string{"versions": {"1", "2", "3"}}
		s.mockVault.EXPECT().ReadData(id, data).Return(&hashicorp.Secret{}, nil)
		s.mockVault.EXPECT().DeleteSecret(id, data).Return(nil)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, false).Return(versions, nil)

		err := s.secretStore.Delete(ctx, id)

		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same NotFound if secret is not found by id ", func() {
		s.mockVault.EXPECT().ReadData(id, gomock.Any()).Return(nil, nil)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, false).Return(versions, nil)

		err := s.secretStore.Delete(ctx, id)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail with same error if delete secret by id fails", func() {
		s.mockVault.EXPECT().ReadData(id, gomock.Any()).Return(&hashicorp.Secret{}, nil)
		s.mockVault.EXPECT().DeleteSecret(id, gomock.Any()).Return(expectedErr)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, false).Return(versions, nil)

		err := s.secretStore.Delete(ctx, id)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestRestore() {
	ctx := context.Background()
	id := "my-restore-secret"
	versions := []string{"1", "2", "3"}

	s.Run("should restore secret by id successfully", func() {
		data := map[string][]string{"versions": {"1", "2", "3"}}
		s.mockVault.EXPECT().RestoreSecret(id, data).Return(nil)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, true).Return(versions, nil)
		err := s.secretStore.Restore(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if restore secret by id fails", func() {
		s.mockVault.EXPECT().RestoreSecret(id, gomock.Any()).Return(expectedErr)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, true).Return(versions, nil)
		err := s.secretStore.Restore(ctx, id)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpSecretStoreTestSuite) TestDestroy() {
	ctx := context.Background()
	id := "my-destroyed-secret"
	versions := []string{"1", "2", "3"}

	s.Run("should destroy secret by id successfully", func() {
		data := map[string][]string{"versions": {"1", "2", "3"}}
		s.mockVault.EXPECT().DestroySecret(id, data).Return(nil)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, true).Return(versions, nil)
		err := s.secretStore.Destroy(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if destroy secret by id fails", func() {
		s.mockVault.EXPECT().DestroySecret(id, gomock.Any()).Return(expectedErr)
		s.mockDB.EXPECT().ListVersions(gomock.Any(), id, true).Return(versions, nil)
		err := s.secretStore.Destroy(ctx, id)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}
