package acceptancetests

import (
	"context"
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type secretsTestSuite struct {
	suite.Suite
	env       *IntegrationEnvironment
	store     stores.SecretStore
	db        database.Secrets
	secretIDs []string
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
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should create a new secret successfully if it already exists in the Vault", func() {
		id := s.newID("my-secret")
		value := "my-secret-value"
		tags := testutils.FakeTags()

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		err = s.db.Delete(ctx, id)
		require.NoError(s.T(), err)
		err = s.db.Purge(ctx, id)
		require.NoError(s.T(), err)

		secret, err = s.store.Set(ctx, id, value, &entities.Attributes{
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
		assert.False(s.T(), secret.Metadata.Disabled)
	})
}

func (s *secretsTestSuite) TestList() {
	ctx := s.env.ctx
	id := s.newID("my-secret-list")
	id2 := s.newID("my-secret-list")
	id3 := s.newID("my-secret-list")
	value := "my-secret-value"

	// 2 with same ID and 1 different
	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id3, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	listLen := 0
	s.Run("should list all secrets ids successfully", func() {
		ids, err := s.store.List(ctx, 0, 0)
		require.NoError(s.T(), err)
		
		listLen = len(ids)
		assert.Contains(s.T(), ids, id)
		assert.Contains(s.T(), ids, id2)
		assert.Contains(s.T(), ids, id3)
	})

	s.Run("should list first secret id successfully", func() {
		ids, err := s.store.List(ctx, 1, uint64(listLen-3))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id})
	})

	s.Run("should list last two secret id successfully", func() {
		ids, err := s.store.List(ctx, 2, uint64(listLen-2))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id2, id3})
	})
}

func (s *secretsTestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-secret-get")
	value := "my-secret-value"

	secret, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)
	version := secret.Metadata.Version

	s.Run("should get secret successfully", func() {
		secret, err := s.store.Get(ctx, id, version)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
		assert.Equal(s.T(), value, secret.Value)
		assert.NotEmpty(s.T(), secret.Metadata.Version)
		assert.NotNil(s.T(), secret.Metadata.CreatedAt)
		assert.NotNil(s.T(), secret.Metadata.UpdatedAt)
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), secret.Metadata.Disabled)
	})

	s.Run("should get latest secret version if no version is specified", func() {
		secret, err := s.store.Get(ctx, id, "")
		require.NoError(s.T(), err)
		assert.Equal(s.T(), version, secret.Metadata.Version)
	})

	s.Run("should fail with NotFound if secret is not found", func() {
		secret, err := s.store.Get(ctx, "nonExistentID", "")

		assert.Nil(s.T(), secret)
		require.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail if version does not exist", func() {
		secret, err := s.store.Get(ctx, id, "41579384e3014e849a2b140463509ea2")

		assert.Nil(s.T(), secret)
		require.Error(s.T(), err)
	})
}

func (s *secretsTestSuite) TestDelete() {
	ctx := s.env.ctx
	id := s.newID("my-secret-delete")
	value := "my-deleted-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)

	s.Run("should delete latest secret successfully", func() {
		err := s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
	})

	s.Run("should fail with NotFound if secret is not found", func() {
		err := s.store.Delete(ctx, "nonExistentID")
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestGetDeleted() {
	ctx := s.env.ctx
	id := fmt.Sprintf("%s-%s", "my-deleted-secret", common.RandString(10))
	value := "my-deleted-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)

	err := s.delete(s.env.ctx, id)
	require.NoError(s.T(), err)

	s.Run("should get deleted secret successfully", func() {
		secret, err := s.store.GetDeleted(ctx, id)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
	})

	s.Run("should fail with NotFound if deleted secret is not found", func() {
		secret, err := s.store.GetDeleted(ctx, "nonExistentID")

		assert.Nil(s.T(), secret)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestRestoredDeletedSecret() {
	ctx := s.env.ctx
	id := s.newID("my-restored-secret")
	value := "my-restored-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)

	err := s.delete(s.env.ctx, id)
	require.NoError(s.T(), err)

	s.Run("should restore deleted secret successfully", func() {
		err := s.store.Restore(ctx, id)
		require.NoError(s.T(), err)
	})

	s.Run("should fail with NotFound if restored secret is not found and not deleted", func() {
		err := s.store.Restore(ctx, "nonExistentID")
		require.NotNil(s.T(), err)
		// require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestListDeleted() {
	ctx := s.env.ctx
	id := s.newID("my-deleted-secret-list")
	id2 := s.newID("my-deleted-secret-list-2")
	id3 := s.newID("my-deleted-secret-list-3")
	value := "my-deleted-secret-value"

	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id3, value, &entities.Attributes{})
	require.NoError(s.T(), err)

	err = s.delete(s.env.ctx, id)
	require.NoError(s.T(), err)

	err = s.delete(s.env.ctx, id2)
	require.NoError(s.T(), err)

	err = s.delete(s.env.ctx, id3)
	require.NoError(s.T(), err)

	listLen := 0
	s.Run("should list all deleted secrets ids successfully", func() {
		ids, err := s.store.ListDeleted(ctx, 0, 0)
		require.NoError(s.T(), err)
		
		listLen = len(ids)
		assert.Contains(s.T(), ids, id)
		assert.Contains(s.T(), ids, id2)
		assert.Contains(s.T(), ids, id3)
	})
	
	s.Run("should list first secret id successfully", func() {
		ids, err := s.store.ListDeleted(ctx, 1, uint64(listLen-3))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id})
	})

	s.Run("should list last two secret id successfully", func() {
		ids, err := s.store.ListDeleted(ctx, 2, uint64(listLen-2))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id2, id3})
	})
}

func (s *secretsTestSuite) newID(name string) string {
	id := fmt.Sprintf("%s-%s", name, common.RandString(10))
	s.secretIDs = append(s.secretIDs, id)

	return id
}

func (s *secretsTestSuite) delete(ctx context.Context, id string) error {
	err := s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	if _, ok := s.store.(*akv.Store); ok {
		maxTries := MaxRetries
		for {
			_, err := s.store.GetDeleted(ctx, id)
			if err == nil {
				break
			}

			if maxTries <= 0 {
				errMsg := "failed to wait for deletion to complete"
				s.env.logger.With("secretID", id).Error(errMsg)
				return fmt.Errorf(errMsg)
			}

			maxTries -= 1
			s.env.logger.Debug("waiting for deletion to complete", "secretID", id)
			time.Sleep(time.Second)
		}
	}

	return nil
}
