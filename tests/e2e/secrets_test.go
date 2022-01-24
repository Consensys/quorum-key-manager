// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type secretsTestSuite struct {
	suite.Suite
	err error
	env *Environment

	storeName string

	deleteQueue  *sync.WaitGroup
	destroyQueue *sync.WaitGroup
}

func TestKeyManagerSecrets(t *testing.T) {
	s := new(secretsTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	if len(s.env.cfg.SecretStores) == 0 {
		t.Error("list of secret stores cannot be empty")
		return
	}

	s.deleteQueue = &sync.WaitGroup{}
	s.destroyQueue = &sync.WaitGroup{}

	for _, storeN := range s.env.cfg.SecretStores {
		s.storeName = storeN
		suite.Run(t, s)
	}
}

func (s *secretsTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.deleteQueue.Wait()
	s.destroyQueue.Wait()
}

func (s *secretsTestSuite) RunT(name string, subtest func()) bool {
	return s.Run(fmt.Sprintf("%s(%s)", name, s.storeName), subtest)
}

func (s *secretsTestSuite) TestSet() {
	secretID := fmt.Sprintf("my-secret-set-%s", common.RandString(10))
	s.RunT("should set a new secret successfully", func() {
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
		require.NoError(s.T(), err)
		defer s.queueToDelete(secret)

		assert.Equal(s.T(), request.Value, secret.Value)
		assert.Equal(s.T(), secretID, secret.ID)
		assert.Equal(s.T(), request.Tags, secret.Tags)
		assert.NotEmpty(s.T(), secret.Version)
		assert.False(s.T(), secret.Disabled)
		assert.NotEmpty(s.T(), secret.CreatedAt)
		assert.NotEmpty(s.T(), secret.UpdatedAt)
	})

	s.RunT("should parse errors successfully", func() {
		secretID := "my-secret-set"
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.env.client.SetSecret(s.env.ctx, "nonExistentStoreName", secretID, request)
		require.Nil(s.T(), secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestGetSecret() {
	secretID := fmt.Sprintf("my-secret-get-%s", common.RandString(10))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)
	time.Sleep(time.Second)

	secret2, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	defer s.queueToDelete(secret)
	defer s.queueToDelete(secret2)

	s.RunT("should get a secret specific version successfully", func() {
		secretRetrieved, err := s.env.client.GetSecret(s.env.ctx, s.storeName, secret.ID, secret.Version)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), request.Value, secretRetrieved.Value)
		assert.Equal(s.T(), secretID, secretRetrieved.ID)
		assert.Equal(s.T(), request.Tags, secretRetrieved.Tags)
		assert.NotEmpty(s.T(), secretRetrieved.Version)
		assert.False(s.T(), secretRetrieved.Disabled)
		assert.NotEmpty(s.T(), secretRetrieved.CreatedAt)
		assert.NotEmpty(s.T(), secretRetrieved.UpdatedAt)
	})

	s.RunT("should get the latest version of a secret successfully", func() {
		secretRetrieved, err := s.env.client.GetSecret(s.env.ctx, s.storeName, secret.ID, "")
		require.NoError(s.T(), err)

		assert.Equal(s.T(), secret2.Version, secretRetrieved.Version)
	})

	s.RunT("should parse errors successfully", func() {
		secret, err := s.env.client.GetSecret(s.env.ctx, s.storeName, secret.ID, "invalidVersion")
		require.Nil(s.T(), secret)

		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestDeleteSecret() {
	secretID := fmt.Sprintf("my-delete-secret-%s", common.RandString(10))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	defer s.queueToDestroy(secret)

	s.RunT("should delete a secret specific version successfully", func() {
		err := s.env.client.DeleteSecret(s.env.ctx, s.storeName, secret.ID)
		assert.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.DeleteSecret(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestGetDeletedSecret() {
	secretID := fmt.Sprintf("my-deleted-secret-%s", common.RandString(10))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	err = s.env.client.DeleteSecret(s.env.ctx, s.storeName, secret.ID)
	require.NoError(s.T(), err)

	defer s.queueToDestroy(secret)

	s.RunT("should get deleted secret successfully", func() {
		secretRetrieved, err := s.env.client.GetDeletedSecret(s.env.ctx, s.storeName, secret.ID)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), secretID, secretRetrieved.ID)
	})

	s.RunT("should parse errors successfully", func() {
		_, err := s.env.client.GetDeletedSecret(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestRestoreDeleted() {
	s.RunT("should restore deleted secret successfully", func() {
		secretID := fmt.Sprintf("my-restore-secret-%s", common.RandString(10))
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
		}

		secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
		require.NoError(s.T(), err)

		err = s.env.client.DeleteSecret(s.env.ctx, s.storeName, secret.ID)
		require.NoError(s.T(), err)
		defer s.queueToDelete(secret)

		// We should retry on status conflict for AKV
		errMsg := fmt.Sprintf("failed to restore secret. {ID: %s}", secret.ID)
		err = retryOn(func() error {
			return s.env.client.RestoreSecret(s.env.ctx, s.storeName, secret.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)
		require.NoError(s.T(), err)

		// We should retry on status conflict for AKV
		errMsg = fmt.Sprintf("failed to get secret. {ID: %s}", secret.ID)
		err = retryOn(func() error {
			_, derr := s.env.client.GetSecret(s.env.ctx, s.storeName, secret.ID, secret.Version)
			return derr
		}, s.T().Logf, errMsg, http.StatusNotFound, MaxRetries)
		require.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.RestoreSecret(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestDestroyDeleted() {
	secretID := fmt.Sprintf("my-destroy-secret-%s", common.RandString(10))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	err = s.env.client.DeleteSecret(s.env.ctx, s.storeName, secret.ID)
	require.NoError(s.T(), err)

	s.RunT("should destroy deleted secret successfully", func() {
		errMsg := fmt.Sprintf("failed to destroy secret {ID: %s}", secret.ID)
		err := retryOn(func() error {
			return s.env.client.DestroySecret(s.env.ctx, s.storeName, secret.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)
		if err != nil {
			httpError, ok := err.(*client.ResponseError)
			require.True(s.T(), ok)
			assert.Equal(s.T(), http.StatusNotImplemented, httpError.StatusCode)
			return
		}
		require.NoError(s.T(), err)

		_, err = s.env.client.GetDeletedSecret(s.env.ctx, s.storeName, secret.ID)
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.env.client.DestroySecret(s.env.ctx, s.storeName, "invalidID")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestList() {
	secretID := fmt.Sprintf("my-secret-list-%s", common.RandString(10))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)
	defer s.queueToDelete(secret)

	s.RunT("should get all secret ids successfully", func() {
		ids, err := s.env.client.ListSecrets(s.env.ctx, s.storeName, 99999, 0)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, secretID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.env.client.ListSecrets(s.env.ctx, "nonExistentStoreName", 0, 0)
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestListDeletedSecrets() {
	s.RunT("should get all deleted secret ids successfully", func() {
		secretID := fmt.Sprintf("my-secret-list-%s", common.RandString(10))
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
		}

		secret, err := s.env.client.SetSecret(s.env.ctx, s.storeName, secretID, request)
		require.NoError(s.T(), err)

		err = s.env.client.DeleteSecret(s.env.ctx, s.storeName, secret.ID)
		require.NoError(s.T(), err)
		defer s.queueToDestroy(secret)

		ids, err := s.env.client.ListDeletedSecrets(s.env.ctx, s.storeName, 99999, 0)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, secretID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.env.client.ListDeletedSecrets(s.env.ctx, "nonExistentStoreName", 0, 0)
		require.Empty(s.T(), ids)

		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) queueToDelete(secretR *types.SecretResponse) {
	s.deleteQueue.Add(1)
	go func() {
		err := s.env.client.DeleteSecret(s.env.ctx, s.storeName, secretR.ID)
		if err != nil {
			s.T().Logf("failed to delete secret {ID: %s}", secretR.ID)
		} else {
			s.queueToDestroy(secretR)
		}
		s.deleteQueue.Done()
	}()
}

func (s *secretsTestSuite) queueToDestroy(secretR *types.SecretResponse) {
	s.destroyQueue.Add(1)
	go func() {
		errMsg := fmt.Sprintf("failed to destroy secret {ID: %s}", secretR.ID)
		err := retryOn(func() error {
			return s.env.client.DestroySecret(s.env.ctx, s.storeName, secretR.ID)
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)

		if err != nil {
			s.T().Logf(errMsg)
		}
		s.destroyQueue.Done()
	}()
}
