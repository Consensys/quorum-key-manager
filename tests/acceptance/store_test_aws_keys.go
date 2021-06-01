// +build acceptance

package acceptancetests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type awsKeysTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *aws.KeyStore
}

func (s *awsKeysTestSuite) TestSet() {
	ctx := s.env.ctx

	s.T().Run("should create a new key successfully", func(t *testing.T) {
		name := "my-key"
		tags := testutils.FakeTags()

		secret, err := s.store.Create(ctx, name, &entities.Algorithm{Type: entities.Ecdsa}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)
		assert.Equal(t, name, secret.ID)

		err = s.store.Destroy(ctx, name)
		require.NoError(s.T(), err)
	})

}
