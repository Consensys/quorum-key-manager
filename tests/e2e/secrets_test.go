// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

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
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *secretsTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})
}

func (s *secretsTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerSecrets(t *testing.T) {
	s := new(secretsTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *secretsTestSuite) TestSet() {
	secretID := fmt.Sprintf("my-secret-set-%d", common.RandInt(1000))
	s.Run("should set a new secret successfully", func() {
		request := &types.SetSecretRequest{
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, secretID, request)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), request.Value, secret.Value)
		assert.Equal(s.T(), secretID, secret.ID)
		assert.Equal(s.T(), request.Tags, secret.Tags)
		assert.Equal(s.T(), "1", secret.Version)
		assert.False(s.T(), secret.Disabled)
		assert.NotEmpty(s.T(), secret.CreatedAt)
		assert.NotEmpty(s.T(), secret.UpdatedAt)
		assert.True(s.T(), secret.ExpireAt.IsZero())
		assert.True(s.T(), secret.DeletedAt.IsZero())
		assert.True(s.T(), secret.DestroyedAt.IsZero())
	})

	s.Run("should parse errors successfully", func() {
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

func (s *secretsTestSuite) TestGet() {
	secretID := fmt.Sprintf("my-secret-get-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, secretID, request)
	require.NoError(s.T(), err)

	secret2, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, secretID, request)
	require.NoError(s.T(), err)

	s.Run("should get a secret specific version successfully", func() {
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, secret.Version)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), request.Value, secretRetrieved.Value)
		assert.Equal(s.T(), secretID, secretRetrieved.ID)
		assert.Equal(s.T(), request.Tags, secretRetrieved.Tags)
		assert.Equal(s.T(), "1", secretRetrieved.Version)
		assert.False(s.T(), secretRetrieved.Disabled)
		assert.NotEmpty(s.T(), secretRetrieved.CreatedAt)
		assert.NotEmpty(s.T(), secretRetrieved.UpdatedAt)
		assert.True(s.T(), secretRetrieved.ExpireAt.IsZero())
		assert.True(s.T(), secretRetrieved.DeletedAt.IsZero())
		assert.True(s.T(), secretRetrieved.DestroyedAt.IsZero())
	})

	s.Run("should get the latest version of a secret successfully", func() {
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, "")
		require.NoError(s.T(), err)

		assert.Equal(s.T(), secret2.Version, secretRetrieved.Version)
	})

	s.Run("should parse errors successfully", func() {
		secret, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, "invalidVersion")
		require.Nil(s.T(), secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 422, httpError.StatusCode)
		assert.Equal(s.T(), "version must be a number", httpError.Message)
	})
}

func (s *secretsTestSuite) TestList() {
	secretID := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		Value: "my-secret-value",
	}

	_, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, secretID, request)
	require.NoError(s.T(), err)

	s.Run("should get all secret ids successfully", func() {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, s.cfg.HashicorpSecretStore)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, secretID)
	})

	s.Run("should parse errors successfully", func() {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}
