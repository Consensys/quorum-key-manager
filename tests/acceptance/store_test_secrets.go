package acceptancetests

import (
	"fmt"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type secretsTestSuite struct {
	suite.Suite
	env       *IntegrationEnvironment
	store     secrets.Store
	secretIDs []string
}

func (s *secretsTestSuite) TearDownSuite() {
	ctx := s.env.ctx

	s.env.logger.Info("Deleting the following secrets", "secrets", s.secretIDs)
	for _, id := range s.secretIDs {
		err := s.store.Delete(ctx, id)
		if err != nil && errors.IsNotSupportedError(err) {
			return
		}
	}

	for _, id := range s.secretIDs {
		maxTries := MAX_RETRIES
		for {
			err := s.store.Destroy(ctx, id)
			if err != nil && !errors.IsStatusConflictError(err) {
				break
			}
			if maxTries <= 0 {
				if err != nil {
					s.env.logger.Info("failed to destroy secret", "secretID", id)
				}
				break
			}

			maxTries -= 1
			waitTime := time.Second * time.Duration(MAX_RETRIES-maxTries)
			s.env.logger.Debug("waiting for deletion to complete", "keyID", id, "waitFor", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}
}

func (s *secretsTestSuite) TestSet() {
	ctx := s.env.ctx

	s.Run("should create a new secret successfully", func() {
		id := s.newID("my-secret")
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
		assert.Equal(s.T(), value, secret.Value)
		assert.Equal(s.T(), tags, secret.Tags)
		assert.NotEmpty(s.T(), secret.Metadata.Version)
		assert.NotNil(s.T(), secret.Metadata.CreatedAt)
		assert.NotNil(s.T(), secret.Metadata.UpdatedAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), secret.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should increase version at each set", func() {
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

		require.NoError(s.T(), err)

		assert.Equal(s.T(), tags1, secret1.Tags)
		assert.Equal(s.T(), value1, secret1.Value)
		assert.Equal(s.T(), tags2, secret2.Tags)
		assert.Equal(s.T(), value2, secret2.Value)
		assert.NotEqual(s.T(), secret1.Metadata.Version, secret2.Metadata.Version)
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

	s.Run("should list all secrets ids successfully", func() {
		ids, err := s.store.List(ctx)

		require.NoError(s.T(), err)
		assert.Contains(s.T(), ids, id)
		assert.Contains(s.T(), ids, id2)
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

	s.Run("should get latest secret successfully if no version is specified", func() {
		secret, err := s.store.Get(ctx, id, "")
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
		assert.Equal(s.T(), value, secret.Value)
		assert.NotEmpty(s.T(), secret.Metadata.Version)
		assert.NotNil(s.T(), secret.Metadata.CreatedAt)
		assert.NotNil(s.T(), secret.Metadata.UpdatedAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), secret.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should get specific secret version", func() {
		secret, err := s.store.Get(ctx, id, version1)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), version1, secret.Metadata.Version)

		secret, err = s.store.Get(ctx, id, version2)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), version2, secret.Metadata.Version)
	})

	s.Run("should fail with NotFound if secret is not found", func() {
		secret, err := s.store.Get(ctx, "inexistentID", "")

		assert.Nil(s.T(), secret)
		require.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail if version does not exist", func() {
		secret, err := s.store.Get(ctx, id, "41579384e3014e849a2b140463509ea2")

		assert.Nil(s.T(), secret)
		require.Error(s.T(), err)
	})
}

func (s *secretsTestSuite) newID(name string) string {
	id := fmt.Sprintf("%s-%d", name, common.RandInt(1000))
	s.secretIDs = append(s.secretIDs, id)

	return id
}
