// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/entities/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"time"

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
	registryName    string
}

func (s *aliasStoreTestSuite) SetupSuite() {
	ctx := context.Background()
	s.user = authtypes.NewWildcardUser()
	s.registryName = "my-registry"
	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	_, err := s.registryService.Create(ctx, s.registryName, []string{s.user.Tenant}, s.user)
	require.NoError(s.T(), err)
}

func (s *aliasStoreTestSuite) fakeAlias() *entities.Alias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return testutils.FakeAlias("JPM-"+randID, "Goldman Sachs-"+randID, entities.AliasKindString, "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=")
}

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
		require.Len(s.T(), registry.Aliases, 2)

		assert.Equal(s.T(), registry.Aliases[0].Key, in.Key)
		assert.Equal(s.T(), registry.Aliases[1].Key, in.Key+"2")
	})

	s.Run("should delete a registry successfully", func() {
		err := s.registryService.Delete(ctx, in.RegistryName, s.user)
		require.NoError(s.T(), err)

		_, err = s.registryService.Get(ctx, in.RegistryName, s.user)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))

		_, err = s.aliasService.Get(ctx, in.RegistryName, in.Key, s.user)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *aliasStoreTestSuite) TestCreateAlias() {
	ctx := context.Background()

	s.Run("should create an alias successfully", func() {
		fakeAlias := s.fakeAlias()
		alias, err := s.aliasService.Create(ctx, s.registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), alias.Key, fakeAlias.Key)
		assert.Equal(s.T(), alias.RegistryName, s.registryName)
		assert.Equal(s.T(), alias.Kind, fakeAlias.Kind)
		assert.Equal(s.T(), alias.Value, fakeAlias.Value)
		assert.NotEmpty(s.T(), alias.CreatedAt)
		assert.NotEmpty(s.T(), alias.UpdatedAt)
		assert.Equal(s.T(), alias.CreatedAt, alias.UpdatedAt)
	})

	s.Run("should fail with NotFoundError if registry does not exist", func() {
		fakeAlias := s.fakeAlias()
		alias, err := s.aliasService.Create(ctx, "inexistent registry", fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
		require.Error(s.T(), err)
		require.Nil(s.T(), alias)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *aliasStoreTestSuite) TestGetAlias() {
	ctx := context.Background()
	fakeAlias := s.fakeAlias()

	createdAlias, err := s.aliasService.Create(ctx, s.registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
	require.NoError(s.T(), err)

	s.Run("should get alias successfully", func() {
		retrievedAlias, err := s.aliasService.Get(ctx, createdAlias.RegistryName, createdAlias.Key, s.user)
		require.NoError(s.T(), err)

		require.Equal(s.T(), retrievedAlias, createdAlias)
	})

	s.Run("should fail with NotFoundError if alias does not exist", func() {
		_, err := s.aliasService.Get(ctx, createdAlias.RegistryName, "inexistentKey", s.user)
		require.Error(s.T(), err)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *aliasStoreTestSuite) TestDeleteAlias() {
	ctx := context.Background()
	fakeAlias := s.fakeAlias()

	createdAlias, err := s.aliasService.Create(ctx, s.registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
	require.NoError(s.T(), err)

	s.Run("just delete alias successfully", func() {
		err = s.aliasService.Delete(ctx, createdAlias.RegistryName, createdAlias.Key, s.user)
		require.NoError(s.T(), err)

		_, err = s.aliasService.Get(ctx, createdAlias.RegistryName, createdAlias.Key, s.user)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should fail with NotFoundError if alias does not exist", func() {
		err := s.aliasService.Delete(ctx, createdAlias.RegistryName, "inexistentAlias", s.user)
		require.Error(s.T(), err)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *aliasStoreTestSuite) TestUpdateAlias() {
	ctx := context.Background()
	fakeAlias := s.fakeAlias()

	createdAlias, err := s.aliasService.Create(ctx, s.registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
	require.NoError(s.T(), err)

	s.Run("should update alias successfully", func() {
		newValue := []interface{}{"SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}
		updatedAlias, err := s.aliasService.Update(
			ctx,
			createdAlias.RegistryName,
			createdAlias.Key,
			entities.AliasKindArray,
			newValue,
			s.user,
		)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), updatedAlias.Key, createdAlias.Key)
		assert.Equal(s.T(), updatedAlias.Kind, entities.AliasKindArray)
		assert.Equal(s.T(), updatedAlias.RegistryName, createdAlias.RegistryName)
		assert.Equal(s.T(), updatedAlias.Value, newValue)
		assert.NotEmpty(s.T(), updatedAlias.CreatedAt)
		assert.NotEmpty(s.T(), updatedAlias.UpdatedAt)
		assert.True(s.T(), updatedAlias.UpdatedAt.After(updatedAlias.CreatedAt))
	})

	s.Run("should fail with NotFoundError if alias does not exist", func() {
		updatedAlias, err := s.aliasService.Update(ctx, createdAlias.RegistryName, "inexistentAlias", createdAlias.Kind, createdAlias.Value, s.user)
		require.Error(s.T(), err)
		require.Nil(s.T(), updatedAlias)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}
