// +build e2e

package e2e

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var encodedApikeySample1 = base64.StdEncoding.EncodeToString([]byte("apikey-sample1"))
var encodedApikeyWrongValue = base64.StdEncoding.EncodeToString([]byte("wrong-apikey"))

type authTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *authTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func (s *authTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

type apiKeyTransport struct{}

// RoundTrip overrides to inject Authorization header with correct APIKey
func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorisation", fmt.Sprintf("Basic %s", encodedApikeySample1))
	return http.DefaultTransport.RoundTrip(req)
}

func TestKeyManagerKeys(t *testing.T) {
	s := new(keysTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *authTestSuite) TestListWithMatchingCert() {
	keyID := "my-key-list"
	request := &types.ImportKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}
	cert, err := tls.LoadX509KeyPair("certs/same.crt", "certs/common.key")
	require.NoError(s.T(), err)
	s.keyManagerClient = client.NewHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	},
	)

	s.Run("should accept request successfully", func() {

		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

type wrongApiKeyTransport struct{}

// RoundTrip overrides to inject Authorization header with wrong APIKey
func (t *wrongApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorisation", fmt.Sprintf("Basic %s", encodedApikeyWrongValue))
	return http.DefaultTransport.RoundTrip(req)
}

func (s *authTestSuite) TestListWithMatchingAPIKey() {

	s.keyManagerClient = client.NewHTTPClient(&http.Client{Transport: &apiKeyTransport{}}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})

	s.Run("should accept request successfully with correct APIKey", func() {

		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

type wrongApiKeyTransport struct{}

// RoundTrip overrides to inject Authorization header with wrong APIKey
func (t *wrongApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorisation", fmt.Sprintf("Basic %s", encodedApikeyWrongValue))
	return http.DefaultTransport.RoundTrip(req)
}

func (s *authTestSuite) TestListWithWrongAPIKey() {

	s.keyManagerClient = client.NewHTTPClient(&http.Client{Transport: &wrongApiKeyTransport{}}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})

	s.Run("should reject request with wrong APIKey", func() {

		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 401, httpError.StatusCode)
	})
}
