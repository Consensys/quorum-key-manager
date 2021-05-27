package acceptancetests

import (
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type secretsTestSuite struct {
	suite.Suite
	env       *IntegrationEnvironment
	store     secrets.Store
	secretIDs []string
}

func (s *secretsTestSuite) TearDownSuite() {
	ctx := s.env.ctx

	s.env.logger.WithField("secrets", s.secretIDs).Info("Deleting the following secrets")
	for _, id := range s.secretIDs {
		_ = s.store.Delete(ctx, id)
	}

	for _, id := range s.secretIDs {
		_ = s.store.Destroy(ctx, id)
	}
}

func (s *secretsTestSuite) TestSet() {
	ctx := s.env.ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		id := s.newID("my-secret")
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
	})

	s.T().Run("should increase version at each set", func(t *testing.T) {
		id := s.newID("my-secret-versioned")
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
	})
}

func (s *secretsTestSuite) TestList() {
	ctx := s.env.ctx
	id := s.newID("my-secret-list")
	id2 := s.newID("my-secret-list")
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
}

func (s *secretsTestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-secret-get")
	value := "my-secret-value"

	// 2 with same ID
	secret1, setErr := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), setErr)
	version1 := secret1.Metadata.Version
	secret2, setErr := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), setErr)
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

	s.T().Run("should fail if version does not exist", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "41579384e3014e849a2b140463509ea2")

		assert.Nil(t, secret)
		require.Error(t, err)
	})
}

func (s *secretsTestSuite) newID(name string) string {
	id := fmt.Sprintf("%s-%d", name, common.RandInt(1000))
	s.secretIDs = append(s.secretIDs, id)

	return id
}
