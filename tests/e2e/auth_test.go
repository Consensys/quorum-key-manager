// +build e2e

package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
	err error
	ctx context.Context

	keyManagerClient *client.HTTPClient
	keyManagerURL    string
	storeName        string
	tlsCACert        string
	tlsCAKey         string
	oidcCAKey        string

	acc    *types.EthAccountResponse
	logger log.Logger
}

func TestAuth(t *testing.T) {
	s := new(authTestSuite)

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
	s.tlsCAKey = cfg.AuthTLSCAKey
	s.tlsCACert = cfg.AuthTLSCACert
	s.oidcCAKey = cfg.AuthOIDCCAKey

	s.logger, s.err = zap.NewLogger(log.NewConfig(log.WarnLevel, log.TextFormat))
	if s.err != nil {
		t.Error(s.err)
		return
	}

	var token string
	token, s.err = generateJWT(s.oidcCAKey, "*:*", "e2e|auth_test")
	if s.err != nil {
		t.Errorf("failed to generate jwt. %s", s.err)
		return
	}
	s.keyManagerClient = client.NewHTTPClient(&http.Client{
		Transport: NewTestHttpTransport(token, "", nil),
	}, &client.Config{
		URL: cfg.KeyManagerURL,
	})

	s.keyManagerURL = cfg.KeyManagerURL
	s.storeName = cfg.EthStores[0]

	suite.Run(t, s)
}

func (s *authTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.acc, s.err = s.keyManagerClient.CreateEthAccount(s.ctx, s.storeName, &types.CreateEthAccountRequest{
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

	_ = s.keyManagerClient.DeleteEthAccount(s.ctx, s.storeName, s.acc.Address.Hex())
	errMsg := fmt.Sprintf("failed to destroy ethAccount {Address: %s}", s.acc.Address.Hex())
	_ = retryOn(func() error {
		return s.keyManagerClient.DestroyEthAccount(s.ctx, s.storeName, s.acc.Address.Hex())
	}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)
}

func (s *authTestSuite) TestAuth_TLS() {
	s.Run("should sign payload successfully", func() {
		clientCert, err := generateClientCert(s.tlsCACert, s.tlsCAKey)
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", "", clientCert),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
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
			URL: s.keyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
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
			URL: s.keyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})
}

func (s *authTestSuite) TestAuth_JWT() {
	s.Run("should sign payload successfully", func() {
		var token string
		token, err := generateJWT(s.oidcCAKey, "*:*", "e2e|auth_test_jwt")
		if s.err != nil {
			s.T().Errorf("failed to generate jwt. %s", s.err)
			return
		}
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		assert.NoError(s.T(), err)
	})

	s.Run("should fail to sign with Status Forbidden", func() {
		var token string
		token, err := generateJWT(s.oidcCAKey, "*:read", "e2e|auth_test_jwt")
		if s.err != nil {
			s.T().Errorf("failed to generate jwt. %s", s.err)
			return
		}
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})

	s.Run("should fail to sign with StatusForbidden", func() {
		var token string
		token, err := generateJWT(s.oidcCAKey, "*:read", "e2e|auth_test_jwt")
		if s.err != nil {
			s.T().Errorf("failed to generate jwt. %s", s.err)
			return
		}
		require.NoError(s.T(), err)

		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err = qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})
}

func (s *authTestSuite) TestAuth_APIKEY() {
	s.Run("should sign payload successfully", func() {
		var apiKey = base64.StdEncoding.EncodeToString([]byte("admin-user"))
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", apiKey, nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		assert.NoError(s.T(), err)
	})

	s.Run("should fail to sign with StatusUnauthorized", func() {
		var apiKey = base64.StdEncoding.EncodeToString([]byte("wrong-apikey"))
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", apiKey, nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusUnauthorized, httpError.StatusCode)
	})

	s.Run("should fail to sign with StatusForbidden", func() {
		var apiKey = base64.StdEncoding.EncodeToString([]byte("guest-user"))
		qkmClient := client.NewHTTPClient(&http.Client{
			Transport: NewTestHttpTransport("", apiKey, nil),
		}, &client.Config{
			URL: s.keyManagerURL,
		})

		_, err := qkmClient.SignMessage(s.ctx, s.storeName, s.acc.Address.Hex(), &types.SignMessageRequest{
			Message: hexutil.MustDecode("0x1234"),
		})
		httpError, ok := err.(*client.ResponseError)
		require.True(s.T(), ok)
		assert.Equal(s.T(), http.StatusForbidden, httpError.StatusCode)
	})
}
