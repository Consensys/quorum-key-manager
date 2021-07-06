// +build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/tests"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type jsonRPCTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
	acc              *types.Eth1AccountResponse
}

func (s *jsonRPCTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})

	var err error
	s.acc, err = s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, &types.CreateEth1AccountRequest{
		KeyID: fmt.Sprintf("test-eth-sign-%d", common.RandInt(1000)),
	})

	if err != nil {
		require.NoError(s.T(), err)
	}
}

func (s *jsonRPCTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	// @TODO validate error once hashicorp support destroy keys
	_ = s.keyManagerClient.DestroyKey(s.ctx, s.cfg.Eth1Store, s.acc.KeyID)
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
	s.Run("should forward call eth_blockNumber and retrieve block number successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_blockNumber")
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		_, err = strconv.ParseUint(result[2:], 16, 64)
		assert.NoError(s.T(), err)
	})
}

func (s *jsonRPCTestSuite) TestEthSign() {
	dataToSign := "0xa2"

	s.Run("should call eth_sign and sign data successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sign", s.acc.Address, dataToSign)
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		// TODO validate signature when importing eth1 accounts is possible
	})

	s.Run("should call eth_sign and fail to sign with an invalid account", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sign", "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852", dataToSign)
		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_sign and fail to sign without an invalid data", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sign", s.acc.Address, "noExpectedHexData")
		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestEthSignTransaction() {
	s.Run("should call eth_signTransaction successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_signTransaction", map[string]interface{}{
			"data":     "0xa2",
			"from":     s.acc.Address,
			"to":       "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"nonce":    "0x0",
			"gas":      "0x989680",
			"gasPrice": "0x10000",
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)
	})

	s.Run("should call eth_signTransaction and fail to sign with an invalid account", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sign", map[string]interface{}{
			"data":     "0xa2",
			"from":     "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852",
			"to":       "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"nonce":    "0x0",
			"gas":      "0x989680",
			"gasPrice": "0x10000",
		})
		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_signTransaction and fail to sign with an invalid args", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sign", "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852", map[string]interface{}{
			"data":  "0xa2",
			"from":  s.acc.Address,
			"to":    "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"nonce": "0x0",
		})

		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestEthSendTransaction() {
	toAddr := "0xd46e8dd67c5d32be8058bb8eb970870f07244567"
	s.Run("should call eth_sendTransaction, successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data": "0xa2",
			"from": s.acc.Address,
			"to":   toAddr,
			"gas":  "0x989680",
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.ctx, s.cfg.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction and fail if an invalid account", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data": "0xa2",
			"from": "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852",
			"to":   toAddr,
			"gas":  "0x989680",
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_sendTransaction and fail if an invalid args", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"from": s.acc.Address,
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestSendPrivTransaction() {
	toAddr := "0xd46e8dd67c5d32be8058bb8eb970870f07244567"

	s.Run("should call eth_sendTransaction, for private Quorum Tx, successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          toAddr,
			"gas":         "0x989680",
			"privateFrom": "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
			"privateFor":  []string{"QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="},
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.ctx, s.cfg.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction and fail if invalid account", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852",
			"to":          toAddr,
			"gas":         "0x989680",
			"privateFrom": "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
			"privateFor":  []string{"QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="},
		})

		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_sendTransaction and fail if invalid args", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          toAddr,
			"gas":         "0x989680",
			"privateFrom": "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})

}

func (s *jsonRPCTestSuite) TestSignEEATransaction() {
	s.Run("should call eea_sendTransaction successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom": "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			"privateFor":  []string{"Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs="},
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.ctx, s.cfg.BesuNodeID, result)
		require.NoError(s.T(), err)
		// Sent to precompiled besu contract
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), "0x000000000000000000000000000000000000007e")
	})

	s.Run("should call eea_sendTransaction and fail if invalid account", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852",
			"to":          "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom": "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			"privateFor":  []string{"Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs="},
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})

	s.Run("should call eea_sendTransaction and fail if invalid args", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom": "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestEthAccounts() {
	s.Run("should call eth_accounts successfully", func() {
		resp, err := s.keyManagerClient.Call(s.ctx, s.cfg.QuorumNodeID, "eth_accounts")
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)
		accs := []string{}
		err = json.Unmarshal(resp.Result.(json.RawMessage), &accs)
		require.NoError(s.T(), err)
		assert.Contains(s.T(), accs, strings.ToLower(s.acc.Address.Hex()))
	})
}

func (s *jsonRPCTestSuite) retrieveTransaction(ctx context.Context, nodeID, txHash string) (*ethtypes.Transaction, error) {
	resp, err := s.keyManagerClient.Call(ctx, nodeID, "eth_getTransactionByHash", txHash)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	var result *ethtypes.Transaction
	err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
	return result, err
}
