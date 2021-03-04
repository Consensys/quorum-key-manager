package secrets

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/infra/vault/mocks"
	"github.com/golang/mock/gomock"
	hashicorp "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type hashicorpSecretStoreTestSuite struct {
	suite.Suite
	mockVault   *mocks.MockHashicorpVaultClient
	mountPoint  string
	secretStore *hashicorpSecretStore
}

func TestHashicorpSecretStore(t *testing.T) {
	s := new(hashicorpSecretStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpSecretStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mountPoint = "secret"
	s.mockVault = mocks.NewMockHashicorpVaultClient(ctrl)

	s.secretStore = New(s.mockVault, s.mountPoint)
}

func (s *hashicorpSecretStoreTestSuite) TestSet() {
	ctx := context.Background()
	id := "my-secret"
	value := "my-value"
	attributes := testutils.FakeAttributes()
	expectedData := map[string]interface{}{
		valueLabel: value,
		tagsLabel:  attributes.Tags,
	}
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			"created_time":  "2018-03-22T02:24:06.945319214Z",
			"deletion_time": "",
			"destroyed":     false,
			"version":       2,
		},
	}

	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")

		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(hashicorpSecret, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)
		assert.False(t, secret.Disabled)
		assert.Equal(t, expectedCreatedAt, secret.CreatedAt)
		assert.Equal(t, attributes.Tags, secret.Tags)
		assert.Equal(t, 2, secret.Version)
		assert.True(t, secret.ExpireAt.IsZero())
		assert.True(t, secret.DeletedAt.IsZero())
		assert.Nil(t, secret.Recovery)
	})

	// TODO: Implement specific error types and check that the function return the right error type
	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(nil, expectedErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.Equal(t, expectedErr, err)
	})

	// TODO: Implement specific error types and check that the function return the right error type
	s.T().Run("should fail with error if it fails to extract metadata", func(t *testing.T) {
		hashSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"created_time": "invalidTime",
				"version":      2,
			},
		}

		s.mockVault.EXPECT().Write(s.mountPoint+"/data/"+id, expectedData).Return(hashSecret, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.Error(t, err)
	})
}
