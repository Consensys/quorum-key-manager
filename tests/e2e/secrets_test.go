// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
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
	id := fmt.Sprintf("my-secret-set-%d", common.RandInt(1000))
	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		request := &types.SetSecretRequest{
			ID:    id,
			Value: "my-secret-value",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		secret, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, request)
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

		secret, err := s.keyManagerClient.SetSecret(s.ctx, "inexistentStoreName", request)
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
	})
}

func (s *secretsTestSuite) TestGet() {
	id := fmt.Sprintf("my-secret-get-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		ID:    id,
		Value: "my-secret-value",
	}

	secret, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, request)
	require.NoError(s.T(), err)

	secret2, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, request)
	require.NoError(s.T(), err)

	s.T().Run("should get a secret specific version successfully", func(t *testing.T) {
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, secret.Version)
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

	s.T().Run("should get the latest version of a secret successfully", func(t *testing.T) {
		secretRetrieved, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, "")
		require.NoError(t, err)

		assert.Equal(t, secret2.Version, secretRetrieved.Version)
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		secret, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpSecretStore, secret.ID, "invalidVersion")
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 422, httpError.StatusCode)
		assert.Equal(t, "version must be a number", httpError.Message)
	})
}

func (s *secretsTestSuite) TestList() {
	id := fmt.Sprintf("my-secret-list-%d", common.RandInt(1000))
	request := &types.SetSecretRequest{
		ID:    id,
		Value: "my-secret-value",
	}

	_, err := s.keyManagerClient.SetSecret(s.ctx, s.cfg.HashicorpSecretStore, request)
	require.NoError(s.T(), err)

	s.T().Run("should get all secret ids successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, s.cfg.HashicorpSecretStore)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(ids), 1)
		assert.Contains(t, ids, request.ID)
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, "inexistentStoreName")
		require.Empty(t, ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
	})
}
