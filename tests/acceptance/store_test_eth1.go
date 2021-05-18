package acceptancetests

import (
	"encoding/base64"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	id := "my-account-set"

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		account, err := s.store.Create(ctx, id, &entities.Attributes{
			Tags: testutils.FakeTags(),
		})

		require.NoError(t, err)
		assert.Equal(t, account.ID, id)
	})
}

func (s *eth1TestSuite) TestSign() {
	ctx := s.env.ctx
	id := "my-account-sign"
	payload := base64.RawURLEncoding.EncodeToString([]byte("my data to sign"))

	account, err := s.store.Import(ctx, id, ecdsaPrivKey, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		signature, err := s.store.Sign(ctx, account.Address, payload)
		require.NoError(t, err)
		assert.Equal(t, "YzQeLIN0Sd43Nbb0QCsVSqChGNAuRaKzEfujnERAJd0523aZyz2KXK93KKh-d4ws3MxAhc8qNG43wYI97Fzi7Q==", signature)

		verified, err := verifySignature(signature, payload, ecdsaPrivKey)
		require.NoError(t, err)
		assert.True(t, verified)
	})
}
