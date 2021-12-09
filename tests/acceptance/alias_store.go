// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
	"math/rand"
	"strconv"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/src/aliases"
)

type aliasStoreTestSuite struct {
	suite.Suite
	env  *IntegrationEnvironment
	srv  aliases.Interactor
	rand *rand.Rand
}

func (s *aliasStoreTestSuite) fakeAlias() entities.Alias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return entities.Alias{
		RegistryName: "JPM-" + randID,
		Key:          "Goldman Sachs-" + randID,
		Value:        entities.AliasValue{Kind: entities.AliasKindString, Value: "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="},
	}
}

func (s *aliasStoreTestSuite) TestCreateAlias() {
	ctx := context.Background()

	s.Run("should create an unique alias without error", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)
	})
}

func (s *aliasStoreTestSuite) TestGetAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.srv.GetAlias(ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		got, err := s.srv.GetAlias(ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &in, got)
	})
}

func (s *aliasStoreTestSuite) TestUpdateAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.srv.UpdateAlias(ctx, in.RegistryName, in)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		updated := in
		updated.Value = entities.AliasValue{Kind: entities.AliasKindArray, Value: []interface{}{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}}

		out, err = s.srv.UpdateAlias(ctx, in.RegistryName, updated)
		require.NoError(s.T(), err)

		got, err := s.srv.GetAlias(ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &updated, got)
	})
}

func (s *aliasStoreTestSuite) TestDeleteAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		err := s.srv.DeleteAlias(ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		err = s.srv.DeleteAlias(ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)

		_, err = s.srv.GetAlias(ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})
}

func (s *aliasStoreTestSuite) TestListAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		als, err := s.srv.ListAliases(ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.Len(s.T(), als, 0)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.srv.CreateAlias(ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		newAlias := in
		newAlias.Key = `Crédit Mutuel`
		newAlias.Value = entities.AliasValue{Kind: entities.AliasKindArray, Value: []interface{}{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="}}
		out, err = s.srv.CreateAlias(ctx, in.RegistryName, newAlias)
		require.NoError(s.T(), err)
		require.Equal(s.T(), newAlias, *out)

		als, err := s.srv.ListAliases(ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), als)
		require.Len(s.T(), als, 2)
		require.Equal(s.T(), als[0].Key, in.Key)
		require.Equal(s.T(), als[0].Value, in.Value)
		require.Equal(s.T(), als[1].Key, newAlias.Key)
		require.Equal(s.T(), als[1].Value, newAlias.Value)
	})
}
