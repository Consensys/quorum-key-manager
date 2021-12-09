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
	err  error
	rand *rand.Rand
	env  *Environment
}

func TestAlias(t *testing.T) {
	env, err := NewEnvironment()
	require.NoError(t, err)

	s := &aliasTestSuite{
		env:  env,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	suite.Run(t, s)
}

type testAlias struct {
	reg    string
	key    string
	val    entities.AliasValue
	newVal entities.AliasValue
}

func (s *aliasTestSuite) fakeAlias() testAlias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testAlias{
		reg:    "JPM-" + randID,
		key:    "GoldmanSachs-" + randID,
		val:    entities.AliasValue{Kind: entities.AliasKindArray, Value: []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}},
		newVal: entities.AliasValue{Kind: entities.AliasKindArray, Value: []interface{}{"ZOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}},
	}
}

func (s *aliasTestSuite) TestFull() {
	fakeAlias := s.fakeAlias()

	s.Run("should create a new alias successfully", func() {
		a, err := s.env.client.CreateAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Kind: fakeAlias.val.Kind, Value: fakeAlias.val.Value})
		s.Require().NoError(err)

		s.Equal(fakeAlias.val.Kind, a.Kind)
		s.Equal(fakeAlias.val.Value, a.Value)
	})

	s.Run("should get the new alias successfully", func() {
		a, err := s.env.client.GetAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.val.Kind, a.Kind)
		s.Equal(fakeAlias.val.Value, a.Value)
	})

	s.Run("should update the new alias with a new value successfully", func() {
		a, err := s.env.client.UpdateAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Kind: fakeAlias.newVal.Kind, Value: fakeAlias.newVal.Value})
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal.Kind, a.Kind)
		s.Equal(fakeAlias.newVal.Value, a.Value)
	})

	s.Run("should get the update alias successfully", func() {
		a, err := s.env.client.GetAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal.Kind, a.Kind)
		s.Equal(fakeAlias.newVal.Value, a.Value)
	})

	s.Run("should list the updated alias successfully", func() {
		as, err := s.env.client.ListAliases(s.env.ctx, fakeAlias.reg)
		s.Require().NoError(err)

		s.Require().Len(as, 1)
		s.Equal(fakeAlias.newVal, as[0].Value)
	})

	s.Run("should delete the updated alias successfully", func() {
		err := s.env.client.DeleteAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)
	})

	s.Run("should fail with not found error if alias is deleted", func() {
		_, err := s.env.client.GetAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestUpdateAlias() {
	fakeAlias := s.fakeAlias()

	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		_, err := s.env.client.UpdateAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Kind: fakeAlias.newVal.Kind, Value: fakeAlias.newVal.Value})
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestGetAlias() {
	fakeAlias := s.fakeAlias()

	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		_, err := s.env.client.GetAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestDeleteAlias() {
	fakeAlias := s.fakeAlias()

	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		err := s.env.client.DeleteAlias(s.env.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestListAliases() {
	fakeAlias := s.fakeAlias()

	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		as, err := s.env.client.ListAliases(s.env.ctx, fakeAlias.reg)
		s.Require().NoError(err)

		s.NotNil(as)
		s.Len(as, 0)
	})
}

func (s *aliasTestSuite) checkErr(err error, status int) {
	httpError, ok := err.(*client.ResponseError)
	s.Require().True(ok)
	s.Equal(status, httpError.StatusCode)
}
