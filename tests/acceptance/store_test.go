// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	memory2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/database/memory"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/local"
	akv2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/akv"
	aws2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/aws"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets/akv"
	aws3 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets/hashicorp"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
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

	// Hashicorp
	logger := log.DefaultLogger().SetComponent("Secrets-Hashicorp")
	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorp.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Secrets-AKV")
	testSuite = new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = akv.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = log.DefaultLogger().SetComponent("Secrets-AWS")
	hashicorpTestSuite := new(awsSecretTestSuite)
	hashicorpTestSuite.env = s.env
	hashicorpTestSuite.store = aws3.New(s.env.awsSecretsClient, logger)
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
	testSuite.store = hashicorp2.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Keys-AKV")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = akv2.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = log.DefaultLogger().SetComponent("Keys-AWS")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = aws2.New(s.env.awsKmsClient, logger)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp
	logger := log.DefaultLogger().SetComponent("Eth1-Hashicorp")
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = local.New(hashicorp2.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), memory2.New(logger), logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = log.DefaultLogger().SetComponent("Eth1-AKV")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = local.New(akv2.New(s.env.akvClient, logger), memory2.New(logger), logger)
	suite.Run(s.T(), testSuite)
}
