// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/entities/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/src/aliases"
)

type aliasStoreTestSuite struct {
	suite.Suite
	env             *IntegrationEnvironment
	aliasService    aliases.Aliases
	registryService aliases.Registries
	user            *authtypes.UserInfo
	rand            *rand.Rand
}

func (s *aliasStoreTestSuite) SetupSuite() {
	s.user = authtypes.NewWildcardUser()
}

func (s *aliasStoreTestSuite) fakeAlias() *entities.Alias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testutils.FakeAlias("JPM-"+randID, "Goldman Sachs-"+randID, entities.AliasKindString, "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=")
}

/*
func (s *aliasStoreTestSuite) TestCreateAlias() {
	ctx := context.Background()

	s.Run("should create an unique alias without error", func() {
		in := s.fakeAlias()
		out, err := s.aliasService.Create(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)
	})
}

func (s *aliasStoreTestSuite) TestGetAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.aliasService.Create(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		got, err := s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &in, got)
	})
}

func (s *aliasStoreTestSuite) TestUpdateAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.aliasService.Update(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.aliasService.Create(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		updated := in

		out, err = s.aliasService.Update(
			ctx,
			in.RegistryName,
			in.Key,
			entities.AliasKindArray, []interface{}{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="},
			s.user,
		)
		require.NoError(s.T(), err)

		got, err := s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &updated, got)
	})
}

func (s *aliasStoreTestSuite) TestDeleteAlias() {
	ctx := context.Background()

	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		err := s.aliasService.Delete(ctx, in.RegistryName, in.Key, s.user)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		out, err := s.aliasService.Create(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)
		require.Equal(s.T(), in, *out)

		err = s.aliasService.Delete(ctx, in.RegistryName, in.Key, s.user)
		require.NoError(s.T(), err)

		_, err = s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}
*/
func (s *aliasStoreTestSuite) TestRegistry() {
	ctx := context.Background()
	in := s.fakeAlias()

	s.Run("should create and get a registry with aliases successfully", func() {
		_, err := s.registryService.Create(ctx, in.RegistryName, []string{}, s.user)
		require.NoError(s.T(), err)

		_, err = s.aliasService.Create(ctx, in.RegistryName, in.Key, in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)

		_, err = s.aliasService.Create(ctx, in.RegistryName, in.Key+"2", in.Kind, in.Value, s.user)
		require.NoError(s.T(), err)

		registry, err := s.registryService.Get(ctx, in.RegistryName, s.user)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), registry.Aliases)

		assert.Len(s.T(), registry.Aliases, 2)
		assert.Equal(s.T(), registry.Aliases[0].Key, in.Key)
		assert.Equal(s.T(), registry.Aliases[1].Key, in.Key+"2")
	})

	s.Run("should delete a registry successfully", func() {
		err := s.registryService.Delete(ctx, in.RegistryName, s.user)
		require.NoError(s.T(), err)

		_, err = s.registryService.Get(ctx, in.RegistryName, s.user)
		require.NoError(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))

		_, err = s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.NoError(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))

	})
}
