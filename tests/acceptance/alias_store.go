// +build acceptance

package acceptancetests

import (
	"math/rand"
	"strconv"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

type aliasStoreTestSuite struct {
	suite.Suite
	env  *IntegrationEnvironment
	srv  aliases.Repository
	rand *rand.Rand
}

func (s *aliasStoreTestSuite) fakeAlias() aliasent.Alias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return aliasent.Alias{
		RegistryName: "JPM-" + randID,
		Key:          "Goldman Sachs-" + randID,
		Value:        []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="},
	}
}

func (s *aliasStoreTestSuite) TestCreateAlias() {
	s.Run("should create an unique alias without error", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)
	})
}

func (s *aliasStoreTestSuite) TestGetAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.srv.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		got, err := s.srv.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &in, got)
	})
}

func (s *aliasStoreTestSuite) TestUpdateAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.srv.UpdateAlias(s.env.ctx, in.RegistryName, in)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		updated := in
		updated.Value = []string{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}

		out, err = s.srv.UpdateAlias(s.env.ctx, in.RegistryName, updated)
		require.NoError(s.T(), err)

		got, err := s.srv.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &updated, got)
	})
}

func (s *aliasStoreTestSuite) TestDeleteAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		err := s.srv.DeleteAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		err = s.srv.DeleteAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)

		_, err = s.srv.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})
}

func (s *aliasStoreTestSuite) TestListAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		als, err := s.srv.ListAliases(s.env.ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.Len(s.T(), als, 0)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		newAlias := in
		newAlias.Key = `Crédit Mutuel`
		newAlias.Value = []string{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="}
		out, err = s.srv.CreateAlias(s.env.ctx, in.RegistryName, newAlias)
		require.NoError(s.T(), err)
		require.Equal(s.T(), newAlias, *out)

		als, err := s.srv.ListAliases(s.env.ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), als)
		require.Len(s.T(), als, 2)
		require.Equal(s.T(), als[0].Key, in.Key)
		require.Equal(s.T(), als[1].Key, newAlias.Key)
	})
}
