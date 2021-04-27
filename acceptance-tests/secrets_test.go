// +build acceptance

package integrationtests

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type secretsTestSuite struct {
	suite.Suite
	env              *IntegrationEnvironment
	err              error
	keyManagerClient *client.HTTPClient
}

func (s *secretsTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *secretsTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerSecrets(t *testing.T) {
	s := new(secretsTestSuite)
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	s.env, err = NewIntegrationEnvironment(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := common.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.env.baseURL,
	})

	suite.Run(t, s)
}

func (s *secretsTestSuite) TestSet() {
	ctx := s.env.ctx

	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		request := &types.SetSecretRequest{
			ID:    "my-secret-set",
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.keyManagerClient.SetSecret(ctx, SecretStoreName, request)
		require.NoError(t, err)

		assert.Equal(t, request.Value, secret.Value)
		assert.Equal(t, request.ID, secret.ID)
		assert.Equal(t, request.Tags, secret.Tags)
		assert.Equal(t, "1", secret.Version)
		assert.False(t, secret.Disabled)
		assert.NotEmpty(t, secret.CreatedAt)
		assert.NotEmpty(t, secret.UpdatedAt)
		assert.True(t, secret.ExpireAt.IsZero())
		assert.True(t, secret.DeletedAt.IsZero())
		assert.True(t, secret.DestroyedAt.IsZero())
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		request := &types.SetSecretRequest{
			ID:    "my-secret-set",
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.keyManagerClient.SetSecret(ctx, "inexistentStoreName", request)
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
		assert.Equal(t, " secret store inexistentStoreName was not found", httpError.Message)
	})
}

func (s *secretsTestSuite) TestGet() {
	ctx := s.env.ctx
	request := &types.SetSecretRequest{
		ID:    "my-secret-get",
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(ctx, SecretStoreName, request)
	require.NoError(s.T(), err)

	s.T().Run("should get a secret successfully", func(t *testing.T) {
		secretRetrieved, err := s.keyManagerClient.GetSecret(ctx, SecretStoreName, secret.ID, secret.Version)
		require.NoError(t, err)

		assert.Equal(t, request.Value, secretRetrieved.Value)
		assert.Equal(t, request.ID, secretRetrieved.ID)
		assert.Equal(t, request.Tags, secretRetrieved.Tags)
		assert.Equal(t, "1", secretRetrieved.Version)
		assert.False(t, secretRetrieved.Disabled)
		assert.NotEmpty(t, secretRetrieved.CreatedAt)
		assert.NotEmpty(t, secretRetrieved.UpdatedAt)
		assert.True(t, secretRetrieved.ExpireAt.IsZero())
		assert.True(t, secretRetrieved.DeletedAt.IsZero())
		assert.True(t, secretRetrieved.DestroyedAt.IsZero())
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		secret, err := s.keyManagerClient.GetSecret(ctx, SecretStoreName, secret.ID, "invalidVersion")
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 422, httpError.StatusCode)
		assert.Equal(t, " version must be a number", httpError.Message)
	})
}

func (s *secretsTestSuite) TestList() {
	ctx := s.env.ctx
	request := &types.SetSecretRequest{
		ID:    "my-secret-list",
		Value: "my-secret-value",
	}

	_, err := s.keyManagerClient.SetSecret(ctx, SecretStoreName, request)
	require.NoError(s.T(), err)

	s.T().Run("should get all secret ids successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(ctx, SecretStoreName)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(ids), 1)
		assert.Contains(t, ids, request.ID)
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(ctx, "inexistentStoreName")
		require.Empty(t, ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
		assert.Equal(t, " secret store inexistentStoreName was not found", httpError.Message)
	})
}
