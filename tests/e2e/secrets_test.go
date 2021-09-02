// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type secretsTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient client.SecretsClient

	storeName string

	deleteQueue  *sync.WaitGroup
	destroyQueue *sync.WaitGroup
}

func TestKeyManagerSecrets(t *testing.T) {
	s := new(secretsTestSuite)
	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	var cfg *tests.Config
	cfg, s.err = tests.NewConfig()
	if s.err != nil {
		t.Error(s.err)
		return
	}

	if len(cfg.SecretStores) == 0 {
		t.Error("list of secret stores cannot be empty")
		return
	}

	s.deleteQueue = &sync.WaitGroup{}
	s.destroyQueue = &sync.WaitGroup{}

	var token string
	token, s.err = generateJWT("./certificates/client.key", "*:*", "e2e|secrets_test")
	if s.err != nil {
		t.Errorf("failed to generate jwt. %s", s.err)
		return
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{
		Transport: NewTestHttpTransport(token, "", nil),
	}, &client.Config{
		URL: cfg.KeyManagerURL,
	})

	for _, storeN := range cfg.SecretStores {
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
	secretID := fmt.Sprintf("my-secret-set-%d", common.RandInt(1000))
	s.RunT("should set a new secret successfully", func() {
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
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

		secret, err := s.keyManagerClient.SetSecret(s.ctx, "inexistentStoreName", secretID, request)
		require.Nil(s.T(), secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestGetSecret() {
	secretID := fmt.Sprintf("my-secret-get-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)
	time.Sleep(time.Second)

	secret2, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	defer s.queueToDelete(secret)
	defer s.queueToDelete(secret2)

	s.RunT("should get a secret specific version successfully", func() {
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.storeName, secret.ID, secret.Version)
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
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.storeName, secret.ID, "")
		require.NoError(s.T(), err)

		assert.Equal(s.T(), secret2.Version, secretRetrieved.Version)
	})

	s.RunT("should parse errors successfully", func() {
		secret, err := s.keyManagerClient.GetSecret(s.ctx, s.storeName, secret.ID, "invalidVersion")
		require.Nil(s.T(), secret)

		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestDeleteSecret() {
	secretID := fmt.Sprintf("my-delete-secret-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	defer s.queueToDestroy(secret)

	s.RunT("should delete a secret specific version successfully", func() {
		err := s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		assert.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, "invalidID", "")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestGetDeletedSecret() {
	secretID := fmt.Sprintf("my-deleted-secret-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	err = s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secret.ID, secret.Version)
	require.NoError(s.T(), err)

	defer s.queueToDestroy(secret)

	s.RunT("should get deleted secret successfully", func() {
		secretRetrieved, err := s.keyManagerClient.GetDeletedSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), secretID, secretRetrieved.ID)
	})

	s.RunT("should parse errors successfully", func() {
		_, err := s.keyManagerClient.GetDeletedSecret(s.ctx, s.storeName, "invalidID", "")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestRestoreDeleted() {
	s.RunT("should restore deleted secret successfully", func() {
		secretID := fmt.Sprintf("my-restore-secret-%d", common.RandInt(1000))
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
		}

		secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
		require.NoError(s.T(), err)

		err = s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		require.NoError(s.T(), err)
		defer s.queueToDelete(secret)

		// We should retry on status conflict for AKV
		errMsg := fmt.Sprintf("failed to restore secret. {ID: %s, version: %s}", secret.ID, secret.Version)
		err = retryOn(func() error {
			return s.keyManagerClient.RestoreSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		}, s.T().Logf, errMsg, http.StatusConflict, MAX_RETRIES)
		require.NoError(s.T(), err)

		// We should retry on status conflict for AKV
		errMsg = fmt.Sprintf("failed to get secret. {ID: %s, version: %s}", secret.ID, secret.Version)
		err = retryOn(func() error {
			_, derr := s.keyManagerClient.GetSecret(s.ctx, s.storeName, secret.ID, secret.Version)
			return derr
		}, s.T().Logf, errMsg, http.StatusNotFound, MAX_RETRIES)
		require.NoError(s.T(), err)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.keyManagerClient.RestoreSecret(s.ctx, s.storeName, "invalidID", "")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestDestroyDeleted() {
	secretID := fmt.Sprintf("my-destroy-secret-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)

	err = s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secret.ID, secret.Version)
	require.NoError(s.T(), err)

	s.RunT("should destroy deleted secret successfully", func() {
		errMsg := fmt.Sprintf("failed to destroy secret {ID: %s, version: %s}", secret.ID, secret.Version)
		err := retryOn(func() error {
			return s.keyManagerClient.DestroySecret(s.ctx, s.storeName, secret.ID, secret.Version)
		}, s.T().Logf, errMsg, http.StatusConflict, MAX_RETRIES)
		require.NoError(s.T(), err)

		_, err = s.keyManagerClient.GetDeletedSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.RunT("should parse errors successfully", func() {
		err := s.keyManagerClient.DestroySecret(s.ctx, s.storeName, "invalidID", "")
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestList() {
	secretID := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
	require.NoError(s.T(), err)
	defer s.queueToDelete(secret)

	s.RunT("should get all secret ids successfully", func() {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, s.storeName)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, secretID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestListDeletedSecrets() {
	s.RunT("should get all deleted secret ids successfully", func() {
		secretID := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
		}

		secret, err := s.keyManagerClient.SetSecret(s.ctx, s.storeName, secretID, request)
		require.NoError(s.T(), err)

		err = s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secret.ID, secret.Version)
		require.NoError(s.T(), err)
		defer s.queueToDestroy(secret)

		ids, err := s.keyManagerClient.ListDeletedSecrets(s.ctx, s.storeName)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, secretID)
	})

	s.RunT("should parse errors successfully", func() {
		ids, err := s.keyManagerClient.ListDeletedSecrets(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) queueToDelete(secretR *types.SecretResponse) {
	s.deleteQueue.Add(1)
	go func() {
		err := s.keyManagerClient.DeleteSecret(s.ctx, s.storeName, secretR.ID, secretR.Version)
		if err != nil {
			s.T().Logf("failed to delete secret {ID: %s, version: %s}", secretR.ID, secretR.Version)
		} else {
			s.queueToDestroy(secretR)
		}
		s.deleteQueue.Done()
	}()
}

func (s *secretsTestSuite) queueToDestroy(secretR *types.SecretResponse) {
	s.destroyQueue.Add(1)
	go func() {
		errMsg := fmt.Sprintf("failed to destroy secret {ID: %s, version: %s}", secretR.ID, secretR.Version)
		err := retryOn(func() error {
			return s.keyManagerClient.DestroySecret(s.ctx, s.storeName, secretR.ID, secretR.Version)
		}, s.T().Logf, errMsg, http.StatusConflict, MAX_RETRIES)

		if err != nil {
			s.T().Logf(errMsg)
		}
		s.destroyQueue.Done()
	}()
}
