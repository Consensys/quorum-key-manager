// +build e2e

package e2e

import (
	"context"
	"crypto/tls"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
)

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

func (s *authTestSuite) TestListWithMatchingCert() {

	cert, err := tls.LoadX509KeyPair("certs/same.crt", "certs/common.key")
	if err != nil {
		s.T().FailNow()
	}
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

func (s *authTestSuite) TestListWithWrongCert() {

	cert, err := tls.LoadX509KeyPair("certs/wrong.crt", "certs/common.key")
	if err != nil {
		s.T().FailNow()
	}
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

	s.Run("should reject request successfully", func() {
		ids, err := s.keyManagerClient.ListKeys(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 401, httpError.StatusCode)
	})
}
