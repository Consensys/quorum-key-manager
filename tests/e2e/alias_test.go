// +build e2e

package e2e

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type aliasTestSuite struct {
	suite.Suite
	err  error
	ctx  context.Context
	rand *rand.Rand

	client *client.HTTPClient
}

func TestAlias(t *testing.T) {
	cfg, err := tests.NewConfig()
	require.NoError(t, err)

	token, err := generateJWT("./certificates/client.key", "*:*", "e2e|keys_test")
	require.NoError(t, err)

	cl := client.NewHTTPClient(
		&http.Client{
			Transport: NewTestHttpTransport(token, "", nil),
		}, &client.Config{
			URL: cfg.KeyManagerURL,
		})
	s := aliasTestSuite{
		client: cl,
		ctx:    context.Background(),
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	suite.Run(t, &s)
}

type testAlias struct {
	reg    types.RegistryName
	key    types.AliasKey
	val    types.AliasValue
	newVal types.AliasValue
}

func (s *aliasTestSuite) fakeAlias() testAlias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testAlias{
		reg:    types.RegistryName("JPM-" + randID),
		key:    types.AliasKey("GoldmanSachs-" + randID),
		val:    `["ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`,
		newVal: `["ZOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`,
	}
}

func (s *aliasTestSuite) TestFull() {
	fakeAlias := s.fakeAlias()
	s.Run("create", func() {
		a, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().NoError(err)

		s.Equal(fakeAlias.val, a.Value)
	})
	s.Run("1st get", func() {
		a, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.val, a.Value)
	})
	s.Run("update", func() {
		a, err := s.client.UpdateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.newVal})
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal, a.Value)
	})
	s.Run("2nd get", func() {
		a, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal, a.Value)
	})
	s.Run("list", func() {
		as, err := s.client.ListAliases(s.ctx, fakeAlias.reg)
		s.Require().NoError(err)

		s.Require().Len(as, 1)
		s.Equal(fakeAlias.newVal, as[0].Value)
	})
	s.Run("delete", func() {
		err := s.client.DeleteAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)
	})
	s.Run("not found get", func() {
		_, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestCreateAlias() {
	fakeAlias := s.fakeAlias()
	s.Run("bad registry format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("bad key format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
}

func (s *aliasTestSuite) checkErr(err error, status int) {
	httpError, ok := err.(*client.ResponseError)
	s.Require().True(ok)
	s.Equal(status, httpError.StatusCode)
}
