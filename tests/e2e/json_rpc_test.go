// +build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type jsonRPCTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *jsonRPCTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})
}

func (s *jsonRPCTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestJSONRpcHTTP(t *testing.T) {
	s := new(jsonRPCTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *jsonRPCTestSuite) TestCallForwarding() {
	s.T().Run("should forward call eth_blockNumber and retrieve block number successfully", func(t *testing.T) {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.NodeID, "eth_blockNumber")
		require.NoError(t, err)
		require.Nil(t, resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(t, err)
		_, err = strconv.ParseUint(result[2:], 16, 64)
		assert.NoError(t, err)
	})
}


func (s *jsonRPCTestSuite) TestEthSign() {
	s.T().Run("should call eth_sign and sign transaction using eth1 store", func(t *testing.T) {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.NodeID, "eth_sign",)
		require.NoError(t, err)
		require.Nil(t, resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(t, err)
		_, err = strconv.ParseUint(result[2:], 16, 64)
		assert.NoError(t, err)
	})
}
