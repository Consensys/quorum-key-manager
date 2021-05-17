package acceptancetests

import (
	"fmt"
	"testing"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/akv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type akvSecretTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *akv.Store
}

func (s *akvSecretTestSuite) TestSet() {
	ctx := s.env.ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		id := fmt.Sprintf("my-secret-%d", common.RandInt(1000))
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, tags, secret.Tags)
		assert.NotEmpty(t, secret.Metadata.Version)
		assert.NotNil(t, secret.Metadata.CreatedAt)
		assert.NotNil(t, secret.Metadata.UpdatedAt)
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
		assert.True(t, secret.Metadata.DestroyedAt.IsZero())
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.False(t, secret.Metadata.Disabled)

		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

	s.T().Run("should increase version at each set", func(t *testing.T) {
		id := fmt.Sprintf("my-secret-versioned-%d", common.RandInt(1000))
		value1 := "my-secret-value1"
		value2 := "my-secret-value2"
		tags1 := testutils.FakeTags()
		tags2 := map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		}

		secret1, err := s.store.Set(ctx, id, value1, &entities.Attributes{
			Tags: tags1,
		})

		secret2, err := s.store.Set(ctx, id, value2, &entities.Attributes{
			Tags: tags2,
		})

		require.NoError(t, err)

		assert.Equal(t, tags1, secret1.Tags)
		assert.Equal(t, value1, secret1.Value)
		assert.Equal(t, tags2, secret2.Tags)
		assert.Equal(t, value2, secret2.Value)
		assert.NotEqual(t, secret1.Metadata.Version, secret2.Metadata.Version)

		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})
}

func (s *akvSecretTestSuite) TestList() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
	id2 := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
	value := "my-secret-value"

	// 2 with same ID and 1 different
	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	s.T().Run("should list all secrets ids successfully", func(t *testing.T) {
		ids, err := s.store.List(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, ids)
	})

	_, err = s.store.Delete(ctx, id)
	require.NoError(s.T(), err)
	_ = s.store.Destroy(ctx, id)

	_, err = s.store.Delete(ctx, id2)
	require.NoError(s.T(), err)
	_ = s.store.Destroy(ctx, id2)
}

func (s *akvSecretTestSuite) TestGet() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-secret-get-%d", common.RandInt(1000))
	value := "my-secret-value"

	// 2 with same ID
	secret1, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	version1 := secret1.Metadata.Version
	secret2, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	version2 := secret2.Metadata.Version

	s.T().Run("should get latest secret successfully if no version is specified", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "")

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
		assert.NotEmpty(t, secret.Metadata.Version)
		assert.NotNil(t, secret.Metadata.CreatedAt)
		assert.NotNil(t, secret.Metadata.UpdatedAt)
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
		assert.True(t, secret.Metadata.DestroyedAt.IsZero())
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.False(t, secret.Metadata.Disabled)
	})

	s.T().Run("should get specific secret version", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, version1)
		require.NoError(t, err)
		assert.Equal(t, version1, secret.Metadata.Version)

		secret, err = s.store.Get(ctx, id, version2)
		require.NoError(t, err)
		assert.Equal(t, version2, secret.Metadata.Version)
	})

	s.T().Run("should fail with NotFound if secret is not found", func(t *testing.T) {
		secret, err := s.store.Get(ctx, "inexistentID", "")

		assert.Nil(t, secret)
		require.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidFormat if version is not formatted correctly", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "invalidVersion")

		assert.Nil(t, secret)
		require.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with NotFound if version does not exist", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "41579384e3014e849a2b140463509ea2")

		assert.Nil(t, secret)
		require.True(t, errors.IsNotFoundError(err))
	})

	_, err = s.store.Delete(ctx, id)
	require.NoError(s.T(), err)
	_ = s.store.Destroy(ctx, id)
}

func (s *akvSecretTestSuite) TestRefresh() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-secret-refresh-%d", common.RandInt(1000))

	s.T().Run("should refresh secret with new expiration date", func(t *testing.T) {
		value1 := "my-secret-value1"
		value2 := "my-secret-value2"

		_, err := s.store.Set(ctx, id, value1, &entities.Attributes{})
		require.NoError(s.T(), err)
		_, err = s.store.Set(ctx, id, value2, &entities.Attributes{})
		require.NoError(s.T(), err)

		err = s.store.Refresh(ctx, id, "", time.Now().Add(time.Hour*24))
		require.NoError(t, err)

		secret, err := s.store.Get(ctx, id, "")

		require.NoError(t, err)
		assert.NotEmpty(t, secret.Metadata.Version)

		assert.True(t, secret.Metadata.ExpireAt.After(time.Now()))

		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

}
