// +build acceptance

package integrationtests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type hashicorpSecretTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *hashicorp.SecretStore
}

func (s *hashicorpSecretTestSuite) TestSet() {
	ctx := s.env.Ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		id := "my-secret"
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
		assert.Equal(t, tags, secret.Tags)
		assert.Equal(t, "1", secret.Metadata.Version)
		assert.NotNil(t, secret.Metadata.CreatedAt)
		assert.NotNil(t, secret.Metadata.UpdatedAt)
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
		assert.True(t, secret.Metadata.DestroyedAt.IsZero())
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.False(t, secret.Metadata.Disabled)
	})

	s.T().Run("should increase version at each set", func(t *testing.T) {
		id := "my-secret-versioned"
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret1, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		secret2, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, "1", secret1.Metadata.Version)
		assert.Equal(t, "2", secret2.Metadata.Version)
	})
}

func (s *hashicorpSecretTestSuite) TestList() {
	ctx := s.env.Ctx
	id := "my-secret1"
	id2 := "my-secret2"
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
		assert.Equal(t, []string{id, id2}, ids)
	})
}
