package acceptancetests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/secrets/hashicorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type hashicorpSecretTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *hashicorp.Store
}

func (s *hashicorpSecretTestSuite) TestSet() {
	ctx := s.env.ctx

	s.Run("should create a new secret successfully", func() {
		id := "my-secret"
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), tags, secret.Tags)
		assert.Equal(s.T(), "1", secret.Metadata.Version)
		assert.NotNil(s.T(), secret.Metadata.CreatedAt)
		assert.NotNil(s.T(), secret.Metadata.UpdatedAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), secret.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should increase version at each set", func() {
		id := "my-secret-versioned"
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

		require.NoError(s.T(), err)

		assert.Equal(s.T(), "1", secret1.Metadata.Version)
		assert.Equal(s.T(), tags1, secret1.Tags)
		assert.Equal(s.T(), value1, secret1.Value)
		assert.Equal(s.T(), "2", secret2.Metadata.Version)
		assert.Equal(s.T(), tags2, secret2.Tags)
		assert.Equal(s.T(), value2, secret2.Value)
	})
}

func (s *hashicorpSecretTestSuite) TestList() {
	ctx := s.env.ctx
	id := "my-secret-list1"
	id2 := "my-secret-list2"
	value := "my-secret-value"

	// 2 with same ID and 1 different
	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	s.Run("should list all secrets ids successfully", func() {
		ids, err := s.store.List(ctx)

		require.NoError(s.T(), err)
		// TODO: Do exact check when Destroy is implemented
		// assert.Equal(s.T(),[]string{id, id2}, ids)
		assert.True(s.T(), len(ids) >= 2)
	})
}

func (s *hashicorpSecretTestSuite) TestGet() {
	ctx := s.env.ctx
	id := "my-secret-get"
	value := "my-secret-value"

	// 2 with same ID
	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	s.Run("should get latest secret successfully if no version is specified", func() {
		secret, err := s.store.Get(ctx, id, "")

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), "2", secret.Metadata.Version)
		assert.NotNil(s.T(), secret.Metadata.CreatedAt)
		assert.NotNil(s.T(), secret.Metadata.UpdatedAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), secret.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should get specific secret version", func() {
		secret1, err := s.store.Get(ctx, id, "1")
		require.NoError(s.T(), err)
		assert.Equal(s.T(), "1", secret1.Metadata.Version)

		secret2, err := s.store.Get(ctx, id, "2")
		require.NoError(s.T(), err)
		assert.Equal(s.T(), "2", secret2.Metadata.Version)

	})

	s.Run("should fail with NotFound if secret is not found", func() {
		secret, err := s.store.Get(ctx, "inexistentID", "")

		assert.Nil(s.T(), secret)
		require.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail with NotFound if version does not exist", func() {
		secret, err := s.store.Get(ctx, id, "3")

		assert.Nil(s.T(), secret)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}
