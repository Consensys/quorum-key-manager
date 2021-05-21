// +build acceptance

package acceptancetests

import (
	"context"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	eth1local "github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1/local"
	hashicorpkey "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/hashicorp"
	"github.com/stretchr/testify/suite"
)

type storeTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *storeTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *storeTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerStore(t *testing.T) {
	s := new(storeTestSuite)
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	s.env, err = NewIntegrationEnvironment(ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := common.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

/*
func (s *storeTestSuite) TestKeyManagerStore_HashicorpSecret() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	logger := log.DefaultLogger().SetComponent("HashicorpSecret")
	store := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)

	testSuite := new(hashicorpSecretTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_AKVSecret() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	logger := log.DefaultLogger().SetComponent("AKVSecret")
	store := akvsecret.New(s.env.akvClient, logger)

	testSuite := new(akvSecretTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_HashicorpKey() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	logger := log.DefaultLogger().SetComponent("HashicorpKey")
	store := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)

	testSuite := new(hashicorpKeyTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_AKVKey() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	logger := log.DefaultLogger().SetComponent("AKVKey")
	store := akvkey.New(s.env.akvClient, logger)

	testSuite := new(akvKeyTestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
*/

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	logger := log.DefaultLogger().SetComponent("Eth1")
	store := eth1local.New(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), logger)

	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
