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

	aliastypes "github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type jsonRPCTestSuite struct {
	suite.Suite
	err             error
	env             *Environment
	acc             *types.EthAccountResponse
	storeName       string
	registryName    string
	alias           string
	ownRegistryName string
	ownAlias        string
	QuorumNodeID    string
	BesuNodeID      string
	GethNodeID      string

	eeaPrivateFromRegistryName          string
	eeaPrivateFromAliasKey              string
	eeaPrivateForRegistryName           string
	eeaPrivateForAliasKey               string
	eeaPrivacyGroupIDStringRegistryName string
	eeaPrivacyGroupIDStringAliasKey     string
	eeaPrivacyGroupIDArrayRegistryName  string
	eeaPrivacyGroupIDArrayAliasKey      string
}

func TestJSONRpcHTTP(t *testing.T) {
	s := new(jsonRPCTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	s.BesuNodeID = s.env.cfg.BesuNodeID
	s.QuorumNodeID = s.env.cfg.QuorumNodeID
	s.GethNodeID = s.env.cfg.GethNodeID
	s.storeName = s.env.cfg.EthStores[0]
	suite.Run(t, s)
}

func (s *jsonRPCTestSuite) SetupSuite() {
	privKey, err := hexutil.Decode("0x56202652fdffd802b7252a456dbd8f3ecc0352bbde76c23b40afe8aebd714e2e")
	if err != nil {
		s.T().Error(err)
	}

	s.acc, err = s.env.client.ImportEthAccount(s.env.ctx, s.storeName, &types.ImportEthAccountRequest{
		KeyID:      fmt.Sprintf("test-eth-sign-%d", common.RandInt(1000)),
		PrivateKey: privKey,
	})
	if err != nil {
		s.T().Error(err)
	}
	s.ownRegistryName = fmt.Sprintf("eth-from-e2e-%s", common.RandString(5))
	s.ownAlias = fmt.Sprintf("eth-from-e2e-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.ownRegistryName, s.ownAlias, aliastypes.AliasRequest{Kind: entities.KindArray, Value: []string{"BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3By="}})
	if err != nil {
		s.T().Error(err)
	}

	s.eeaPrivateFromRegistryName = fmt.Sprintf("eea-from-e2e-%s", common.RandString(5))
	s.eeaPrivateFromAliasKey = fmt.Sprintf("eea-from-e2e-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.eeaPrivateFromRegistryName, s.eeaPrivateFromAliasKey, aliastypes.AliasRequest{Kind: entities.KindString, Value: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="})
	if err != nil {
		s.T().Error(err)
	}

	s.eeaPrivateForRegistryName = fmt.Sprintf("eea-for-e2e-%s", common.RandString(5))
	s.eeaPrivateForAliasKey = fmt.Sprintf("eea-for-e2e-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.eeaPrivateForRegistryName, s.eeaPrivateForAliasKey, aliastypes.AliasRequest{Kind: entities.KindArray, Value: []interface{}{"Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs="}})
	if err != nil {
		s.T().Error(err)
	}

	s.eeaPrivacyGroupIDStringRegistryName = fmt.Sprintf("eea-groupIDString-e2e-%s", common.RandString(5))
	s.eeaPrivacyGroupIDStringAliasKey = fmt.Sprintf("eea-groupIDString-e2e-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.eeaPrivacyGroupIDStringRegistryName, s.eeaPrivacyGroupIDStringAliasKey, aliastypes.AliasRequest{Kind: entities.KindString, Value: "// TODO: CHANGE"})
	if err != nil {
		s.T().Error(err)
	}

	s.eeaPrivacyGroupIDArrayRegistryName = fmt.Sprintf("eea-groupIDArray-e2e-%s", common.RandString(5))
	s.eeaPrivacyGroupIDArrayAliasKey = fmt.Sprintf("eea-groupIDArray-e2e-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.eeaPrivacyGroupIDArrayRegistryName, s.eeaPrivacyGroupIDArrayAliasKey, aliastypes.AliasRequest{Kind: entities.KindArray, Value: []string{"// TODO: CHANGE"}})
	if err != nil {
		s.T().Error(err)
	}

	s.registryName = fmt.Sprintf("e2e-%s", common.RandString(5))
	s.alias = fmt.Sprintf("Group-A-%s", common.RandString(5))
	_, err = s.env.client.CreateAlias(s.env.ctx, s.registryName, s.alias, aliastypes.AliasRequest{Kind: entities.KindArray, Value: []string{"QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="}})
	if err != nil {
		s.T().Error(err)
	}
}

func (s *jsonRPCTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	err := s.env.client.DeleteEthAccount(s.env.ctx, s.storeName, s.acc.Address.Hex())
	if err != nil {
		s.T().Error(err)
	}

	err = retryOn(func() error {
		return s.env.client.DestroyEthAccount(s.env.ctx, s.storeName, s.acc.Address.Hex())
	}, s.T().Logf, fmt.Sprintf("failed to destroy ethAccount {Address: %s}", s.acc.Address.Hex()), http.StatusConflict, MaxRetries)
	if err != nil {
		s.T().Error(err)
	}

	err = s.env.client.DeleteAlias(s.env.ctx, s.registryName, s.alias)
	if err != nil {
		s.T().Error(err)
	}
}

func (s *jsonRPCTestSuite) TestCallForwarding() {
	s.Run("should forward call eth_blockNumber and retrieve block number successfully", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_blockNumber")
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
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sign", s.acc.Address, dataToSign)
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		// TODO validate signature when importing eth accounts is possible
	})

	s.Run("should call eth_sign and fail to sign with an invalid account", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sign", "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852", dataToSign)
		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_sign and fail to sign without an invalid data", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sign", s.acc.Address, "noExpectedHexData")
		require.NoError(s.T(), err)
		require.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestEthSignTransaction() {
	s.Run("should call eth_signTransaction successfully; legacy tx", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_signTransaction", map[string]interface{}{
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

	s.Run("should call eth_signTransaction successfully; dynamic fee tx", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_signTransaction", map[string]interface{}{
			"data":                 "0xa2",
			"from":                 s.acc.Address,
			"to":                   "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"nonce":                "0x1",
			"gas":                  "0x989680",
			"maxFeePerGas":         "0x10000",
			"maxPriorityFeePerGas": "0x1000",
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)
	})

	s.Run("should call eth_signTransaction and fail to sign with an invalid account", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sign", map[string]interface{}{
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
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sign", "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852", map[string]interface{}{
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

	s.Run("should call eth_sendTransaction successfully: legacy tx", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.GethNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":     "0xa2",
			"from":     s.acc.Address,
			"to":       toAddr,
			"gasPrice": "0x3b9aca07",
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.GethNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction successfully: dynamic fee tx", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.GethNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":  "0xa2",
			"from":  s.acc.Address,
			"to":    toAddr,
			"value": "0x1",
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.GethNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction and fail if an invalid account", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data": "0xa2",
			"from": "0xeE4ec3235F4b08ADC64f539BaC598c5E67BdA852",
			"to":   toAddr,
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})

	s.Run("should call eth_sendTransaction and fail if an invalid args", func() {
		// some args are missing
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"from": s.acc.Address,
			"gas":  "0x989680",
		})

		require.NoError(s.T(), err)
		assert.Error(s.T(), resp.Error)
	})
}

func (s *jsonRPCTestSuite) TestSendPrivTransaction() {
	toAddr := "0xd46e8dd67c5d32be8058bb8eb970870f07244567"

	s.Run("should call eth_sendTransaction, for private Quorum Tx, successfully", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
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
		tx, err := s.retrieveTransaction(s.env.ctx, s.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction, for private Quorum Tx, with a privateFrom alias successfully", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          toAddr,
			"gas":         "0x989680",
			"privateFrom": fmt.Sprintf("{{%s:%s}}", s.ownRegistryName, s.ownAlias),
			"privateFor":  []string{fmt.Sprintf("{{%s:%s}}", s.registryName, s.alias)},
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction, for private Quorum Tx, with an privateFor alias successfully", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          toAddr,
			"gas":         "0x989680",
			"privateFrom": "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
			"privateFor":  []string{fmt.Sprintf("{{%s:%s}}", s.registryName, s.alias)},
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction, for private Quorum Tx, with a privacyGroupID alias successfully", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
			"data":           "0xa2",
			"from":           s.acc.Address,
			"to":             toAddr,
			"gas":            "0x989680",
			"privateFrom":    "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
			"privacyGroupID": fmt.Sprintf("{{%s:%s}}", s.registryName, s.alias),
		})
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.QuorumNodeID, result)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), toAddr)
	})

	s.Run("should call eth_sendTransaction and fail if invalid account", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
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
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_sendTransaction", map[string]interface{}{
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
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
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
		tx, err := s.retrieveTransaction(s.env.ctx, s.BesuNodeID, result)
		require.NoError(s.T(), err)
		// Sent to precompiled besu contract
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), "0x000000000000000000000000000000000000007e")
	})

	s.Run("should call eea_sendTransaction successfully, with alias privateFrom", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom": fmt.Sprintf("{{%s:%s}}", s.eeaPrivateFromRegistryName, s.eeaPrivateFromAliasKey),
			"privateFor":  []string{"Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs="},
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.BesuNodeID, result)
		require.NoError(s.T(), err)
		// Sent to precompiled besu contract
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), "0x000000000000000000000000000000000000007e")
	})

	s.Run("should call eea_sendTransaction successfully, with alias privateFor", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":        "0xa2",
			"from":        s.acc.Address,
			"to":          "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom": "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			"privateFor":  []string{fmt.Sprintf("{{%s:%s}}", s.eeaPrivateForRegistryName, s.eeaPrivateForAliasKey)},
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.BesuNodeID, result)
		require.NoError(s.T(), err)
		// Sent to precompiled besu contract
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), "0x000000000000000000000000000000000000007e")
	})

	s.Run("should call eea_sendTransaction successfully, with alias privacyGroupID as array", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
			"data":           "0xa2",
			"from":           s.acc.Address,
			"to":             "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
			"privateFrom":    "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
			"privacyGroupId": fmt.Sprintf("{{%s:%s}}", s.eeaPrivateForRegistryName, s.eeaPrivateForAliasKey),
		})

		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)

		var result string
		err = json.Unmarshal(resp.Result.(json.RawMessage), &result)
		assert.NoError(s.T(), err)
		tx, err := s.retrieveTransaction(s.env.ctx, s.BesuNodeID, result)
		require.NoError(s.T(), err)
		// Sent to precompiled besu contract
		assert.Equal(s.T(), strings.ToLower(tx.To().String()), "0x000000000000000000000000000000000000007e")
	})

	s.Run("should call eea_sendTransaction and fail if invalid account", func() {
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
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
		resp, err := s.env.client.Call(s.env.ctx, s.BesuNodeID, "eea_sendTransaction", map[string]interface{}{
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
		resp, err := s.env.client.Call(s.env.ctx, s.QuorumNodeID, "eth_accounts")
		require.NoError(s.T(), err)
		require.Nil(s.T(), resp.Error)
		accs := []string{}
		err = json.Unmarshal(resp.Result.(json.RawMessage), &accs)
		require.NoError(s.T(), err)
		assert.Contains(s.T(), accs, strings.ToLower(s.acc.Address.Hex()))
	})
}

func (s *jsonRPCTestSuite) retrieveTransaction(ctx context.Context, nodeID, txHash string) (*ethtypes.Transaction, error) {
	resp, err := s.env.client.Call(ctx, nodeID, "eth_getTransactionByHash", txHash)
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
