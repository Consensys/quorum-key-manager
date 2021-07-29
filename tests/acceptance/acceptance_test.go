// +build acceptance

package acceptancetests

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/postgres"
	eth1 "github.com/consensys/quorum-key-manager/src/stores/store/eth1/local"
	akvkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
	awskey "github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/local"
	akvsecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	awssecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/aws"
	hashicorpsecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/common"
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

	// Hashicorp
	logger := s.env.logger.WithComponent("Secrets-Hashicorp")
	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Secrets-AKV")
	testSuite = new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = akvsecret.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)
	
	// AWS
	logger = s.env.logger.WithComponent("Secrets-AWS")
	testSuite = new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.store = awssecret.New(s.env.awsSecretsClient, logger)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Keys() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	// Hashicorp
	logger := s.env.logger.WithComponent("Keys-Hashicorp")
	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Keys-AKV")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = akvkey.New(s.env.akvClient, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = s.env.logger.WithComponent("Keys-AWS")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = awskey.New(s.env.awsKmsClient, logger)
	suite.Run(s.T(), testSuite)

	// Local
	logger = s.env.logger.WithComponent("Keys-Local")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	db := postgres.New(logger, s.env.postgresClient)
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	testSuite.store = local.New(hashicorpSecretStore, db, logger)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Eth1-DB"), s.env.postgresClient)

	// Hashicorp
	logger := s.env.logger.WithComponent("Eth1-Hashicorp")
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.New(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Eth1-AKV")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.New(akvkey.New(s.env.akvClient, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AWS
	logger = s.env.logger.WithComponent("Eth1-AWS")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.New(awskey.New(s.env.awsKmsClient, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// Local
	logger = s.env.logger.WithComponent("Eth1-Local")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	testSuite.store = eth1.New(local.New(hashicorpSecretStore, db, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)
}

// Please keep this function to clean the keys
/*
func cleanKeys(ctx context.Context, store keys.Store) error {
	keyIDs, err := store.List(ctx)
	if err != nil {
		return err
	}

	for len(keyIDs) != 0 {
		for _, id := range keyIDs {
			err = store.Delete(ctx, id)
			if err != nil {
				return err
			}
		}

		keyIDs, err = store.List(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
*/
