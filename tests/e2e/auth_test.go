// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
	err       error
	env       *Environment
	storeName string
	acc       *types.EthAccountResponse
}

func TestAuth(t *testing.T) {
	s := new(authTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	s.storeName = s.env.cfg.EthStores[0]

	suite.Run(t, s)
}

func (s *authTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.acc, s.err = s.env.client.CreateEthAccount(s.env.ctx, s.storeName, &types.CreateEthAccountRequest{
		KeyID: fmt.Sprintf("e2e-auth-test-%d", common.RandInt(1000)),
	})

	if s.err != nil {
		s.T().Error(s.err)
	}
}

func (s *authTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	_ = s.env.client.DeleteEthAccount(s.env.ctx, s.storeName, s.acc.Address.Hex())
	errMsg := fmt.Sprintf("failed to destroy ethAccount {Address: %s}", s.acc.Address.Hex())
	_ = retryOn(func() error {
		return s.env.client.DestroyEthAccount(s.env.ctx, s.storeName, s.acc.Address.Hex())
	}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)
}

func (s *authTestSuite) TestAuth_TLS() {
	s.Run("should sign payload successfully", func() {
		clientCert, err := generateClientCert(s.env.cfg.AuthTLSCert, s.env.cfg.AuthTLSKey)
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "", clientCert),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		assert.NoError(s.T(), err)
	})

	s.Run("should fail to sign with Status Forbidden if the certificate does not contain auth info", func() {
		clientCert, err := generateClientCert("./certificates/client_no_auth.crt", "./certificates/client_no_auth.key")
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "", clientCert),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})

	s.Run("should fail to sign with StatusForbidden", func() {
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})
}

func (s *authTestSuite) TestAuth_JWT() {
	s.Run("should sign payload successfully", func() {
		token, err := getJWT(s.env.cfg.AuthOIDCTokenURL, s.env.cfg.AuthOIDCClientID, s.env.cfg.AuthOIDCClientSecret, "https://quorum-key-manager.consensys.net/admin")
		if s.err != nil {
			s.T().Errorf("failed to generate jwt. %s", s.err)
			return
		}
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		assert.NoError(s.T(), err)
	})

	s.Run("should fail to sign with status Forbidden if token is not authorized", func() {
		token, err := getJWT(s.env.cfg.AuthOIDCTokenURL, s.env.cfg.AuthOIDCClientID, s.env.cfg.AuthOIDCClientSecret, "https://quorum-key-manager.consensys.net")
		if s.err != nil {
			s.T().Errorf("failed to generate jwt. %s", s.err)
			return
		}
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})

	s.Run("should fail to sign with Status Unauthorized if token is invalid", func() {
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("invalidToken", "", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusUnauthorized, httpError.StatusCode)
	})
}

func (s *authTestSuite) TestAuth_APIKEY() {
	s.Run("should sign payload successfully", func() {
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "admin-user", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		assert.NoError(s.T(), err)
	})

	s.Run("should fail to sign with StatusUnauthorized", func() {
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "wrong-apikey", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusUnauthorized, httpError.StatusCode)
	})

	s.Run("should fail to sign with StatusForbidden", func() {
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "guest-user", nil),
		}, &client.Config{
			URL: s.env.cfg.KeyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.env.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})
}
