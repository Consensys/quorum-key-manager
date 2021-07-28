package acceptancetests

import (
	"context"
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/aws"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
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
		maxTries := MaxRetries
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
			waitTime := time.Second * time.Duration(MaxRetries-maxTries)
			s.env.logger.Debug("waiting for deletion to complete", "secretID", id, "waitFor", waitTime.Seconds())
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

	secret, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)
	version := secret.Metadata.Version

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
		secret, err := s.store.Get(ctx, id, version)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), version, secret.Metadata.Version)
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
		err := s.store.Delete(ctx, "inexistentID")
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestGetDeleted() {
	// Skip not supported secret store types
	if _, ok := s.store.(*hashicorp.Store); ok {
        return
    } 
    
	if _, ok := s.store.(*aws.Store); ok {
        return
    }

	ctx := s.env.ctx
	id := s.newID("my-deleted-secret")
	value := "my-deleted-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)
	deleteErr := s.store.Delete(ctx, id)
	require.NoError(s.T(), deleteErr)
	
	// Wait time for key to transition to Deleted state
	if _, ok := s.store.(*akv.Store); ok {
		err := s.waitDeletedStatus(s.env.ctx, id)
		require.NoError(s.T(), err)
    }

	s.Run("should get deleted secret successfully", func() {
		secret, err := s.store.GetDeleted(ctx, id)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, secret.ID)
	})

	s.Run("should fail with NotFound if deleted secret is not found", func() {
		secret, err := s.store.GetDeleted(ctx, "inexistentID")

		assert.Nil(s.T(), secret)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}


func (s *secretsTestSuite) TestRestoredDeletedSecret() {
	// Skip not supported secret store types
	if _, ok := s.store.(*hashicorp.Store); ok {
        return
    } 

	ctx := s.env.ctx
	id := s.newID("my-restored-secret")
	value := "my-restored-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)

	deleteErr := s.store.Delete(ctx, id)
	require.NoError(s.T(), deleteErr)
	
	// Wait time for key to transition to Deleted state
	if _, ok := s.store.(*akv.Store); ok {
		err := s.waitDeletedStatus(s.env.ctx, id)
		require.NoError(s.T(), err)
    }
	
	s.Run("should restore deleted secret successfully", func() {
		err := s.store.Undelete(ctx, id)
		require.NoError(s.T(), err)
	})

	s.Run("should fail with NotFound if deleted secret is not found", func() {
		err := s.store.Undelete(ctx, "inexistentID")
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestDestroyDeletedSecret() {
	// Skip not supported secret store types
	if _, ok := s.store.(*hashicorp.Store); ok {
        return
    } 
    
	ctx := s.env.ctx
	id := s.newID("my-destroy-secret")
	value := "my-destroy-secret-value"

	_, setErr := s.store.Set(ctx, id, value, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), setErr)

	deleteErr := s.store.Delete(ctx, id)
	require.NoError(s.T(), deleteErr)
	
	// Wait time for key to transition to Deleted state
	if _, ok := s.store.(*akv.Store); ok {
		err := s.waitDeletedStatus(s.env.ctx, id)
		require.NoError(s.T(), err)
    }
	
	s.Run("should destroy deleted secret successfully", func() {
		err := s.store.Destroy(ctx, id)
		require.NoError(s.T(), err)
		
		secret, err := s.store.GetDeleted(ctx, id)
		require.Nil(s.T(), secret)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail with NotFound if deleted secret is not found", func() {
		err := s.store.Undelete(ctx, "inexistentID")
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *secretsTestSuite) TestListDeleted() {
	// Skip not supported secret store types
	if _, ok := s.store.(*hashicorp.Store); ok {
        return
    } 
    
	if _, ok := s.store.(*aws.Store); ok {
        return
    }

	ctx := s.env.ctx
	id := s.newID("my-deleted-secret-list")
	id2 := s.newID("my-deleted-secret-list-2")
	value := "my-deleted-secret-value"

	_, err := s.store.Set(ctx, id, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	_, err = s.store.Set(ctx, id2, value, &entities.Attributes{})
	require.NoError(s.T(), err)
	
	err = s.store.Delete(ctx, id)
	require.NoError(s.T(), err)
	err = s.store.Delete(ctx, id2)
	require.NoError(s.T(), err)
	
	// Wait time for key to transition to Deleted state
	if _, ok := s.store.(*akv.Store); ok {
		err := s.waitDeletedStatus(s.env.ctx, id)
		require.NoError(s.T(), err)
    }
	

	s.Run("should list all deleted secrets ids successfully", func() {
		ids, err := s.store.ListDeleted(ctx)

		require.NoError(s.T(), err)
		assert.Contains(s.T(), ids, id)
		assert.Contains(s.T(), ids, id2)
	})
}

func (s *secretsTestSuite) newID(name string) string {
	id := fmt.Sprintf("%s-%s", name, common.RandString(10))
	s.secretIDs = append(s.secretIDs, id)

	return id
}

func (s *secretsTestSuite) waitDeletedStatus(ctx context.Context, id string) error {
	maxTries := MaxRetries
	for {
		_, err := s.store.GetDeleted(ctx, id)
		if err == nil  {
			break
		} else if !errors.IsNotFoundError(err) {
			s.env.logger.Error("failed to get deleted secret", "secretID", id)
			return err
		}

		if maxTries <= 0 {
			if err != nil {
				errMsg := "failed to wait for deletion to complete"
				s.env.logger.Error(errMsg, "secretID", id)
				return fmt.Errorf(errMsg)
			}
			break
		}

		maxTries -= 1
		waitTime := time.Second * time.Duration(MaxRetries-maxTries)
		s.env.logger.Debug("waiting for deletion to complete", "secretID", id, "waitFor", waitTime.Seconds())
		time.Sleep(waitTime)
	}

	return nil
}
