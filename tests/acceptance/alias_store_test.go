// +build acceptance

package acceptancetests

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/pkg/common"
	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/models"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/store/postgres"
)

func TestAliasStore(t *testing.T) {
	s := new(aliasStoreTestSuite)
	ctx, cancel := context.WithCancel(context.Background())

	sig := common.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	var err error
	s.env, err = NewIntegrationEnvironment(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}
	suite.Run(t, s)
}

type aliasStoreTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store aliasstore.Store
	rand  *rand.Rand
}

func (s *aliasStoreTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}
	s.env.logger.Info("setup test suite has completed")

	s.store = aliaspg.New(s.env.postgresClient)
	randSrc := rand.NewSource(time.Now().UnixNano())
	s.rand = rand.New(randSrc)
}

func (s *aliasStoreTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
}

func (s *aliasStoreTestSuite) fakeAlias() aliasmodels.Alias {
	randInt := s.rand.Intn(1 << 32)
	randID := strconv.Itoa(randInt)
	return aliasmodels.Alias{
		RegistryName: aliasmodels.RegistryName("JPM-" + randID),
		Key:          aliasmodels.AliasKey("Goldman Sachs-" + randID),
		Value:        `["ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`,
	}
}

func (s *aliasStoreTestSuite) TestCreateAlias() {
	s.Run("should create an unique alias without error", func() {
		in := s.fakeAlias()
		err := s.store.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)
	})
}

func (s *aliasStoreTestSuite) TestGetAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		_, err := s.store.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		err := s.store.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)

		got, err := s.store.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &in, got)
	})
}

func (s *aliasStoreTestSuite) TestUpdateAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		err := s.store.UpdateAlias(s.env.ctx, in.RegistryName, in)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		err := s.store.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)

		updated := in
		updated.Value = `["SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","3T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`

		err = s.store.UpdateAlias(s.env.ctx, in.RegistryName, updated)
		require.NoError(s.T(), err)

		got, err := s.store.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)
		require.Equal(s.T(), &updated, got)
	})
}

func (s *aliasStoreTestSuite) TestDeleteAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		err := s.store.DeleteAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		err := s.store.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)

		err = s.store.DeleteAlias(s.env.ctx, in.RegistryName, in.Key)
		require.NoError(s.T(), err)

		_, err = s.store.GetAlias(s.env.ctx, in.RegistryName, in.Key)
		require.Error(s.T(), err)
	})
}

func (s *aliasStoreTestSuite) TestListAlias() {
	s.Run("non existing alias", func() {
		in := s.fakeAlias()
		als, err := s.store.ListAliases(s.env.ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.Len(s.T(), als, 0)
	})

	s.Run("just created alias", func() {
		in := s.fakeAlias()
		err := s.store.CreateAlias(s.env.ctx, in.RegistryName, in)
		require.NoError(s.T(), err)

		newAlias := in
		newAlias.Key = `Crédit Mutuel`
		newAlias.Value = `[ SOAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc= ]`
		err = s.store.CreateAlias(s.env.ctx, in.RegistryName, newAlias)
		require.NoError(s.T(), err)

		als, err := s.store.ListAliases(s.env.ctx, in.RegistryName)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), als)
		require.Len(s.T(), als, 2)
		require.Equal(s.T(), als[0].Key, in.Key)
		require.Equal(s.T(), als[1].Key, newAlias.Key)
	})
}
