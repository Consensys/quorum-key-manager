// +build acceptance

package acceptancetests

import (
	entities2 "github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	testutils2 "github.com/consensysquorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys/aws"
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
		tags := testutils2.FakeTags()

		secret, err := s.store.Create(ctx, name, &entities2.Algorithm{Type: entities2.Ecdsa}, &entities2.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)
		assert.Equal(t, name, secret.ID)

		err = s.store.Destroy(ctx, name)
		require.NoError(s.T(), err)
	})

}
