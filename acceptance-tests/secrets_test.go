// +build acceptance

package integrationtests

import (
	"fmt"
	client "github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/types/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

// secretsTestSuite is a test suite for Key Manager secrets
type secretsTestSuite struct {
	suite.Suite
	baseURL string
	client  client.KeyManagerClient
	env     *IntegrationEnvironment
}

func (s *secretsTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL)
	s.client = client.NewHTTPClient(&http.Client{}, conf)
}

func (s *secretsTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		secretRequest := testutils.FakeCreateSecretRequest()

		secret, err := s.client.CreateSecret(ctx, secretRequest)

		assert.NoError(t, err)
		fmt.Println(secret)
	})
}
