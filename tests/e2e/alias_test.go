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
	reg    string
	key    string
	val    []string
	newVal []string
}

func (s *aliasTestSuite) fakeAlias() testAlias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testAlias{
		reg:    "JPM-" + randID,
		key:    "GoldmanSachs-" + randID,
		val:    []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="},
		newVal: []string{"ZOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="},
	}
}

func (s *aliasTestSuite) TestFull() {
	fakeAlias := s.fakeAlias()
	s.Run("should create a new alias successfully", func() {
		a, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().NoError(err)

		s.Equal(fakeAlias.val, a.Value)
	})
	s.Run("should get the new alias successfully", func() {
		a, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.val, a.Value)
	})
	s.Run("should update the new alias with a new value successfully", func() {
		a, err := s.client.UpdateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.newVal})
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal, a.Value)
	})
	s.Run("should get the update alias successfully", func() {
		a, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)

		s.Equal(fakeAlias.newVal, a.Value)
	})
	s.Run("should list the updated alias successfully", func() {
		as, err := s.client.ListAliases(s.ctx, fakeAlias.reg)
		s.Require().NoError(err)

		s.Require().Len(as, 1)
		s.Equal(fakeAlias.newVal, as[0].Value)
	})
	s.Run("should delete the updated alias successfully", func() {
		err := s.client.DeleteAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().NoError(err)
	})
	s.Run("should fail with not found error if alias is deleted", func() {
		_, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestCreateAlias() {
	fakeAlias := s.fakeAlias()
	s.Run("should fail with bad request if registry has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should not fail if registry use chars that need to be URL encoded", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad registry"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with bad request if key has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should not fail if key has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		_, err := s.client.CreateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.val})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
}

func (s *aliasTestSuite) TestUpdateAlias() {
	fakeAlias := s.fakeAlias()
	s.Run("should fail with bad request if registry has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		_, err := s.client.UpdateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.newVal})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with bad request if key has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		_, err := s.client.UpdateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.newVal})
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		_, err := s.client.UpdateAlias(s.ctx, fakeAlias.reg, fakeAlias.key, types.AliasRequest{Value: fakeAlias.newVal})
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestGetAlias() {
	fakeAlias := s.fakeAlias()
	s.Run("should fail with bad request if registry has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		_, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with bad request if key has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		_, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		_, err := s.client.GetAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestDeleteAlias() {
	fakeAlias := s.fakeAlias()
	s.Run("should fail with bad request if registry has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		err := s.client.DeleteAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with bad request if key has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "bad@key"
		err := s.client.DeleteAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		err := s.client.DeleteAlias(s.ctx, fakeAlias.reg, fakeAlias.key)
		s.Require().Error(err)

		s.checkErr(err, http.StatusNotFound)
	})
}

func (s *aliasTestSuite) TestListAliases() {
	fakeAlias := s.fakeAlias()
	s.Run("should fail with bad request if registry has a bad format", func() {
		fakeAlias := fakeAlias
		fakeAlias.reg = "bad@registry"
		_, err := s.client.ListAliases(s.ctx, fakeAlias.reg)
		s.Require().Error(err)

		s.checkErr(err, http.StatusBadRequest)
	})
	s.Run("should fail with not found if key does not exist", func() {
		fakeAlias := fakeAlias
		fakeAlias.key = "notfound-key"
		as, err := s.client.ListAliases(s.ctx, fakeAlias.reg)
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
