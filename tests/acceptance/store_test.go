// +build acceptance

package acceptancetests

import (
	"context"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/database/memory"
	eth1 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/local"
	akvkey "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/akv"
	hashicorpkey "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/hashicorp"
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
func (s *storeTestSuite) TestKeyManagerStore_Secrets() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp
	logger := log.DefaultLogger().SetComponent("Secrets-Hashicorp")
	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Secrets-AKV")
	testSuite = new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = akvsecret.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = log.DefaultLogger().SetComponent("Secrets-AWS")
	hashicorpTestSuite := new(awsSecretTestSuite)
	hashicorpTestSuite.env = s.env
	hashicorpTestSuite.store = aws.New(s.env.awsVaultClient, logger)
	suite.Run(s.T(), hashicorpTestSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Keys() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp
	logger := log.DefaultLogger().SetComponent("Keys-Hashicorp")
	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Keys-AKV")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = akvkey.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)
}
*/
func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp
	logger := log.DefaultLogger().SetComponent("Eth1-Hashicorp")
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.New(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), memory.New(logger), logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Eth1-AKV")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.New(akvkey.New(s.env.akvClient, logger), memory.New(logger), logger)
	suite.Run(s.T(), testSuite)
}
