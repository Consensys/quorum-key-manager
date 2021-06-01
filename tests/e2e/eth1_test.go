// +build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types/testutils"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// const ecdsaPrivKeyHex = ""

type eth1TestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *eth1TestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})
}

func (s *eth1TestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerEth1(t *testing.T) {
	s := new(eth1TestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *eth1TestSuite) TestCreate() {
	s.Run("should create a new account successfully", func() {
		request := testutils.FakeCreateEth1AccountRequest()

		acc, err := s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.Equal(s.T(), request.ID, acc.ID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
		assert.True(s.T(), acc.ExpireAt.IsZero())
		assert.True(s.T(), acc.DeletedAt.IsZero())
		assert.True(s.T(), acc.DestroyedAt.IsZero())
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeCreateEth1AccountRequest()

		key, err := s.keyManagerClient.CreateEth1Account(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}
