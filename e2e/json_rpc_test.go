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
}

func (s *jsonRPCTestSuite) SetupSuite() {
	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: keyManagerURL,
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

	suite.Run(t, s)
}

func (s *jsonRPCTestSuite) TestRequestForwarding() {
	s.T().Run("should forward call eth_blockNumber", func(t *testing.T) {
		resp, err := s.keyManagerClient.Call(s.ctx, nodeID, "eth_blockNumber")
		require.NoError(t, err)

		assert.Empty(t, resp.Error)
		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(t, err)
		n, err := strconv.ParseUint(result[2:], 16, 64)
		assert.NoError(t, err)
		assert.NotEmpty(t, n)
	})
}
