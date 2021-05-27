// +build acceptance

package acceptancetests

import (
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type awsSecretTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *aws.SecretStore
}

func (s *awsSecretTestSuite) TestSet() {
	ctx := s.env.ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		name := "my-secret"
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, name, value, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)
		assert.Equal(t, name, secret.ID)

		err = s.store.Destroy(ctx, name)
		require.NoError(s.T(), err)
	})

	s.T().Run("should increase version at each set", func(t *testing.T) {
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

		require.NoError(t, err)

		assert.Equal(t, tags1, secret1.Tags)
		assert.Equal(t, value1, secret1.Value)
		assert.Equal(t, tags2, secret2.Tags)
		assert.Equal(t, value2, secret2.Value)
		assert.NotEqual(t, secret1.Metadata.Version, secret2.Metadata.Version)

		err = s.store.Destroy(ctx, id)
		require.NoError(s.T(), err)
	})
}

func (s *awsSecretTestSuite) TestList() {
	ctx := s.env.ctx
	id1 := "my-secret-list1"
	id2 := "my-secret-list2"
	value := "my-secret-value"

	// 2 with same ID and 1 different
	_, err := s.store.Set(ctx, id1, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	s.T().Run("should list all secrets ids successfully", func(t *testing.T) {
		ids, err := s.store.List(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, ids)
		assert.Contains(t, ids, id1)
		assert.Contains(t, ids, id2)
	})

	err = s.store.Destroy(ctx, id1)
	require.NoError(s.T(), err)
	err = s.store.Destroy(ctx, id2)
	require.NoError(s.T(), err)

	// 100 with random IDs ID
	randomIDs := make([]string, 100)
	randomValues := make([]string, 100)

	for i := 0; i < len(randomIDs); i++ {
		randomIDs[i] = fmt.Sprintf("randomID%d", common.RandInt(100000))
		randomValues[i] = fmt.Sprintf("randomValues%d", common.RandInt(100000))
		s.store.Set(ctx, randomIDs[i], randomValues[i], &entities.Attributes{})
	}

	s.T().Run("should list all secrets ids successfully", func(t *testing.T) {
		ids, err := s.store.List(ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, ids)
		assert.NotNil(t, ids)
	})

	for i := 0; i < len(randomIDs); i++ {
		err = s.store.Destroy(ctx, randomIDs[i])
		require.NoError(s.T(), err)
	}
}

func (s *awsSecretTestSuite) TestGet() {
	ctx := s.env.ctx
	id := "my-secret-get"
	value1 := "my-secret-value1"
	value2 := "my-secret-value2"

	// 2 with same ID
	secret1, err := s.store.Set(ctx, id, value1, &entities.Attributes{})
	require.NoError(s.T(), err)
	version1 := secret1.Metadata.Version
	secret2, err := s.store.Set(ctx, id, value2, &entities.Attributes{})
	require.NoError(s.T(), err)
	version2 := secret2.Metadata.Version

	s.T().Run("should get latest secret successfully if no version is specified", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "")

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value2, secret.Value)
		assert.NotEmpty(t, secret.Metadata.Version)
		assert.NotNil(t, secret.Metadata.CreatedAt)
		assert.NotNil(t, secret.Metadata.UpdatedAt)
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
		assert.True(t, secret.Metadata.DestroyedAt.IsZero())
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.False(t, secret.Metadata.Disabled)
	})

	s.T().Run("should get specific secret version", func(t *testing.T) {
		readSec1, err := s.store.Get(ctx, id, version1)
		require.NoError(t, err)

		readSec2, err := s.store.Get(ctx, id, version2)
		require.NoError(t, err)
		expectedVersion2 := readSec2.Metadata.Version
		assert.Equal(t, version2, expectedVersion2)
		assert.False(t, assert.ObjectsAreEqualValues(readSec2, readSec1))
	})

	s.T().Run("should fail with NotFound if secret is not found", func(t *testing.T) {
		secret, err := s.store.Get(ctx, "nonexistentID", "")

		assert.Nil(t, secret)
		require.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with NotFound if version does not exist", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "41579384e3014e849a2b140463509ea2")

		assert.Nil(t, secret)
		require.True(t, errors.IsNotFoundError(err))
	})

	err = s.store.Destroy(ctx, id)
	require.NoError(s.T(), err)
}

func (s *awsSecretTestSuite) TestDeleteAndDestroy() {

	ctx := s.env.ctx
	id := "my-secret-destroy"
	value := "my-secret-value"

	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	s.T().Run("should get secret successfully before destroyed", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "")

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
	})

	s.T().Run("should raise a not found error when deleted", func(t *testing.T) {
		err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_, err := s.store.Get(ctx, id, "")

		require.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should Undelete existing secret", func(t *testing.T) {
		err = s.store.Undelete(ctx, id)
		require.NoError(s.T(), err)
	})

	s.T().Run("should find Secret again when Undeleted", func(t *testing.T) {
		secret, err := s.store.Get(ctx, id, "")

		require.NoError(t, err)

		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
	})

	err = s.store.Destroy(ctx, id)
	require.NoError(s.T(), err)

	s.T().Run("should raise a not found error when destroyed", func(t *testing.T) {
		_, err := s.store.Get(ctx, id, "")

		require.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should list Zero secrets ids", func(t *testing.T) {
		ids, err := s.store.List(ctx)

		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}
