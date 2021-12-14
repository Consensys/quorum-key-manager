// +build e2e

package e2e

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/stretchr/testify/suite"
)

type aliasTestSuite struct {
	suite.Suite
	err   error
	rand  *rand.Rand
	env   *Environment
	alias testAlias
}

func TestAlias(t *testing.T) {
	env, err := NewEnvironment()
	require.NoError(t, err)

	s := &aliasTestSuite{
		env:  env,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	s.alias = s.fakeAlias()

	suite.Run(t, s)
}

func (s *aliasTestSuite) TestFull() {
	s.Run("should create a new registry successfully", func() {
		registry, err := s.env.client.CreateRegistry(s.env.ctx, s.alias.reg, &types.CreateRegistryRequest{
			AllowedTenants: []string{"tenant1"},
		})
		s.Require().NoError(err)

		s.Equal(s.alias.reg, registry.Name)
		s.NotEmpty(registry.UpdatedAt)
		s.NotEmpty(registry.CreatedAt)
		s.True(registry.CreatedAt.Equal(registry.UpdatedAt))
	})

	s.Run("should create a new alias successfully", func() {
		a, err := s.env.client.CreateAlias(s.env.ctx, s.alias.reg, s.alias.key, &types.AliasRequest{Kind: s.alias.kind, Value: s.alias.val})
		s.Require().NoError(err)

		s.Equal(s.alias.reg, a.Registry)
		s.Equal(s.alias.kind, a.Kind)
		s.Equal(s.alias.key, a.Key)
		s.Equal(s.alias.val, a.Value)
		s.NotEmpty(a.UpdatedAt)
		s.NotEmpty(a.CreatedAt)
		s.True(a.CreatedAt.Equal(a.UpdatedAt))
	})

	s.Run("should get the new alias successfully", func() {
		a, err := s.env.client.GetAlias(s.env.ctx, s.alias.reg, s.alias.key)
		s.Require().NoError(err)

		s.Equal(s.alias.reg, a.Registry)
		s.Equal(s.alias.kind, a.Kind)
		s.Equal(s.alias.key, a.Key)
		s.Equal(s.alias.val, a.Value)
		s.NotEmpty(a.UpdatedAt)
		s.NotEmpty(a.CreatedAt)
		s.True(a.CreatedAt.Equal(a.UpdatedAt))
	})

	s.Run("should update the new alias with a new value successfully", func() {
		a, err := s.env.client.UpdateAlias(s.env.ctx, s.alias.reg, s.alias.key, &types.AliasRequest{Kind: s.alias.newKind, Value: s.alias.newVal})
		s.Require().NoError(err)

		s.Equal(s.alias.newKind, a.Kind)
		s.Equal(s.alias.newVal, a.Value)
	})

	s.Run("should list the aliases of the registry successfully", func() {
		registry, err := s.env.client.GetRegistry(s.env.ctx, s.alias.reg)
		s.Require().NoError(err)

		s.Require().Len(registry.Aliases, 1)
		s.Equal(s.alias.newVal, registry.Aliases[0].Value)
	})

	s.Run("should delete the registry successfully", func() {
		err := s.env.client.DeleteRegistry(s.env.ctx, s.alias.reg)
		s.Require().NoError(err)
	})

	s.Run("should fail with not found error if registry is deleted", func() {
		_, err := s.env.client.GetAlias(s.env.ctx, s.alias.reg, s.alias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestUpdateAlias() {
	s.Run("should fail with not found if key does not exist", func() {
		_, err := s.env.client.UpdateAlias(s.env.ctx, s.alias.reg, "notfound-key", &types.AliasRequest{Kind: s.alias.newKind, Value: s.alias.newVal})
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestGetAlias() {
	s.Run("should fail with not found if key does not exist", func() {
		_, err := s.env.client.GetAlias(s.env.ctx, s.alias.reg, "notfound-key")
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestDeleteAlias() {
	s.Run("should fail with not found if key does not exist", func() {
		err := s.env.client.DeleteAlias(s.env.ctx, s.alias.reg, "notfound-key")
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

type testAlias struct {
	reg     string
	key     string
	kind    string
	val     interface{}
	newKind string
	newVal  interface{}
}

func (s *aliasTestSuite) fakeAlias() testAlias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testAlias{
		reg:     "JPM-" + randID,
		key:     "GoldmanSachs-" + randID,
		kind:    entities.AliasKindArray,
		val:     []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="},
		newKind: entities.AliasKindString,
		newVal:  "ZOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=",
	}
}

func (s *aliasTestSuite) checkErr(err error, status int) {
	httpError, ok := err.(*client.ResponseError)
	s.Require().True(ok)
	s.Equal(status, httpError.StatusCode)
}
