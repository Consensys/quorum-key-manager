package acceptancetests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type eth1TestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store eth1.Store
}

func (s *eth1TestSuite) TestSet() {
	ctx := s.env.ctx
	id := "my-account"

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		account, err := s.store.Create(ctx, id, &entities.Attributes{
			Tags: testutils.FakeTags(),
		})

		assert.NoError(t, err)
		assert.Equal(t, account.ID, id)
	})
}
