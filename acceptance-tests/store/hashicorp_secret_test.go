// +build acceptance

package store

import (
	integrationtests "github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type hashicorpSecretTestSuite struct {
	suite.Suite
	env   *integrationtests.IntegrationEnvironment
	store *hashicorp.SecretStore
}

func (s *hashicorpSecretTestSuite) TestCreate() {
	ctx := s.env.Ctx

	s.T().Run("should create a new secret successfully", func(t *testing.T) {
		id := "my-secret"
		value := "my-secret-value"

		secret, err := s.store.Set(ctx, id, value, &entities.Attributes{
			Tags: testutils.FakeTags(),
		})

		assert.NoError(t, err)
		assert.Equal(t, id, secret.ID)
		assert.Equal(t, value, secret.Value)
	})
}
