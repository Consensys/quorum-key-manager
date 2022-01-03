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
	env              *IntegrationEnvironment
	aliasService     aliases.Aliases
	registryService  aliases.Registries
	user             *authtypes.UserInfo
	userUnauthorized *authtypes.UserInfo
	rand             *rand.Rand
	registryName     string
}

func (s *aliasStoreTestSuite) SetupSuite() {
	ctx := context.Background()
	s.user = authtypes.NewWildcardUser()
	s.user.Tenant = "tenantAllowed"
	s.userUnauthorized = authtypes.NewAnonymousUser()
	s.userUnauthorized.Tenant = "tenantUnauthorized"
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

	s.Run("should fail to create a registry with ForbiddenError if missing permission", func() {
		registry, err := s.registryService.Create(ctx, in.RegistryName, []string{}, s.userUnauthorized)
		s.Error(err)
		s.Nil(registry)

		s.True(errors.IsForbiddenError(err))
	})

	s.Run("should create and get a registry with aliases successfully", func() {
		_, err := s.registryService.Create(ctx, in.RegistryName, []string{s.user.Tenant}, s.user)
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

	s.Run("should fail to delete a registry with ForbiddenError if missing permission", func() {
		err := s.registryService.Delete(ctx, in.RegistryName, s.userUnauthorized)
		s.Error(err)

		s.True(errors.IsForbiddenError(err))
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

	s.Run("should fail with ForbiddenError if missing permission", func() {
		fakeAlias := s.fakeAlias()
		registry, err := s.aliasService.Create(ctx, s.registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.userUnauthorized)
		s.Error(err)
		s.Nil(registry)

		s.True(errors.IsForbiddenError(err))
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

	s.Run("should fail with ForbiddenError if missing permission", func() {
		registry, err := s.aliasService.Get(ctx, createdAlias.RegistryName, createdAlias.Key, s.userUnauthorized)
		s.Error(err)
		s.Nil(registry)

		s.True(errors.IsForbiddenError(err))
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

	s.Run("should fail with ForbiddenError if missing permission", func() {
		err := s.aliasService.Delete(ctx, createdAlias.RegistryName, createdAlias.Key, s.userUnauthorized)
		s.Error(err)

		s.True(errors.IsForbiddenError(err))
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

	s.Run("should fail with ForbiddenError if missing permission", func() {
		registry, err := s.aliasService.Update(ctx, createdAlias.RegistryName, createdAlias.Key, createdAlias.Kind, createdAlias.Value, s.userUnauthorized)
		s.Error(err)
		s.Nil(registry)

		s.True(errors.IsForbiddenError(err))
	})
}

func (s *aliasStoreTestSuite) TestAccess() {
	ctx := context.Background()
	registryName := "my-restricted-registry"
	fakeAlias := s.fakeAlias()
	userNoAccess := authtypes.NewAnonymousUser()
	userNoAccess.Tenant = "tenantUnauthorized"
	userNoAccess.Permissions = authtypes.ListWildcardPermission("*:aliases")

	restrictedRegistry, err := s.registryService.Create(ctx, registryName, []string{s.user.Tenant}, s.user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), restrictedRegistry.AllowedTenants, []string{s.user.Tenant})

	s.Run("should fail to get registry with NotFoundError if not allowed ", func() {
		registry, err := s.registryService.Get(ctx, registryName, userNoAccess)
		require.Error(s.T(), err)
		require.Nil(s.T(), registry)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should get registry successfully if allowed", func() {
		registry, err := s.registryService.Get(ctx, registryName, s.user)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), registry.Name, registryName)
	})

	s.Run("should fail to insert an alias in a registry with NotFoundError if not allowed ", func() {
		alias, err := s.aliasService.Create(ctx, registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, userNoAccess)
		require.Error(s.T(), err)
		require.Nil(s.T(), alias)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should insert an alias in a registry successfully if allowed ", func() {
		_, err := s.aliasService.Create(ctx, registryName, fakeAlias.Key, fakeAlias.Kind, fakeAlias.Value, s.user)
		require.NoError(s.T(), err)
	})

	s.Run("should fail to update an alias in a registry with NotFoundError if not allowed ", func() {
		alias, err := s.aliasService.Update(ctx, registryName, fakeAlias.Key, entities.AliasKindString, "my-new-alias", userNoAccess)
		require.Error(s.T(), err)
		require.Nil(s.T(), alias)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should update an alias in a registry successfully if allowed ", func() {
		alias, err := s.aliasService.Update(ctx, registryName, fakeAlias.Key, entities.AliasKindString, "my-new-alias", s.user)
		require.NoError(s.T(), err)

		value, _ := alias.String()
		assert.Equal(s.T(), "my-new-alias", value)
	})

	s.Run("should fail to delete an alias in a registry with NotFoundError if not allowed ", func() {
		err := s.aliasService.Delete(ctx, registryName, fakeAlias.Key, userNoAccess)
		require.Error(s.T(), err)

		assert.True(s.T(), errors.IsNotFoundError(err))
	})

	s.Run("should delete an alias in a registry successfully if allowed ", func() {
		err := s.aliasService.Delete(ctx, registryName, fakeAlias.Key, s.user)
		require.NoError(s.T(), err)
	})
}
