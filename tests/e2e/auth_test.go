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
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ecdsaPrivKey, _ = base64.StdEncoding.DecodeString("2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=")
var eddsaPrivKey, _ = base64.StdEncoding.DecodeString("X9Yz/5+O42+eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZGGbioek5qYuzJzTNZpTHrVjjFk7iFe3FYwfpxZyNPxtIaFB5gb9VP9IcHZewwNZly821re7RkmB8pGdjywygPH")

type keysTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *keysTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})
}

func (s *keysTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
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

func (s *keysTestSuite) TestListWithMatchingCert() {
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
	if err != nil {
		s.T().FailNow()
	}
	s.keyManagerClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	s.Run("should accept request successfully", func() {

		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestListWithWrongCert() {
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
	cert, err := tls.LoadX509KeyPair("certs/wrong.crt", "certs/common.key")
	if err != nil {
		s.T().FailNow()
	}
	s.keyManagerClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	s.Run("should reject request successfully", func() {
		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 401, httpError.StatusCode)
	})
}
