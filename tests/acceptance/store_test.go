// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/database/memory"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	eth1local "github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1/local"
	akvkey "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/akv"
	hashicorpkey "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/hashicorp"
	akvsecret "github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/akv"
	hashicorpsecret "github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
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

func (s *storeTestSuite) TestKeyManagerStore_Secrets() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp tests
	logger := log.DefaultLogger().SetComponent("Secrets-Hashicorp")
	hashicorpStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)

	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpStore
	suite.Run(s.T(), testSuite)

	// AKV test suite
	logger = log.DefaultLogger().SetComponent("Secrets-AKV")
	akvStore := akvsecret.New(s.env.akvClient, logger)

	testSuite = new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = akvStore
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Keys() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp tests
	logger := log.DefaultLogger().SetComponent("Keys-Hashicorp")
	hashicorpStore := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)

	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpStore
	suite.Run(s.T(), testSuite)

	// AKV test suite
	logger = log.DefaultLogger().SetComponent("Keys-AKV")
	akvStore := akvkey.New(s.env.akvClient, logger)

	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = akvStore
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp tests
	logger := log.DefaultLogger().SetComponent("Eth1-Hashicorp")
	testSuite := new(eth1TestSuite)
	hashicorpKeyStore := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	eth1AccountsDB := memory.New(logger)
	store := eth1local.New(hashicorpKeyStore, eth1AccountsDB, logger)

	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)

	// AKV test suite
	logger = log.DefaultLogger().SetComponent("Eth1-AKV")
	testSuite = new(eth1TestSuite)
	akvKeyStore := akvkey.New(s.env.akvClient, logger)
	eth1AccountsDB = memory.New(logger)
	store = eth1local.New(akvKeyStore, eth1AccountsDB, logger)

	testSuite.env = s.env
	testSuite.store = store
	suite.Run(s.T(), testSuite)
}
