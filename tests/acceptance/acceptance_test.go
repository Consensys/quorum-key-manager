// +build acceptance

package acceptancetests

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/eth1"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	akvkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
	awskey "github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/local"
	hashicorpsecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
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

func (s *storeTestSuite) TestKeyManager_Keys() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Keys-DB"), s.env.postgresClient)

	// Hashicorp
	logger := s.env.logger.WithComponent("Keys-Hashicorp")
	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = keys.NewConnector(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), db, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Keys-AKV")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = keys.NewConnector(akvkey.New(s.env.akvClient, logger), db, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = s.env.logger.WithComponent("Keys-AWS")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.store = keys.NewConnector(awskey.New(s.env.awsKmsClient, db.Keys(), logger), db, logger)
	suite.Run(s.T(), testSuite)

	// Local
	logger = s.env.logger.WithComponent("Keys-Local")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	testSuite.store = keys.NewConnector(local.New(hashicorpSecretStore, logger), db, logger)
	suite.Run(s.T(), testSuite)
}
*/

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Eth1-DB"), s.env.postgresClient)

	// Hashicorp
	storeName := "Eth1-Hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	hashicorpStore := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(hashicorpStore, db.ETH1Accounts(storeName), logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AKV
	storeName = "Eth1-AKV"
	logger = s.env.logger.WithComponent(storeName)
	akvStore := akvkey.New(s.env.akvClient, logger)
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(akvStore, db.ETH1Accounts(storeName), logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AWS
	storeName = "Eth1-AWS"
	logger = s.env.logger.WithComponent(storeName)
	awsStore := awskey.New(s.env.awsKmsClient, logger)
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(awsStore, db.ETH1Accounts(storeName), logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// Local
	storeName = "Eth1-Local-Hashicorp"
	logger = s.env.logger.WithComponent(storeName)
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	localStore := local.New(hashicorpSecretStore, logger)
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(localStore, db.ETH1Accounts(storeName), logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

}

// Please keep this function to clean the keys
func cleanAKVKeys(ctx context.Context, akvClient akv.Client) error {
	kItems, err := akvClient.GetKeys(ctx, 0)
	if err != nil {
		return err
	}

	for len(kItems) != 0 {
		for _, kItem := range kItems {
			_, err = akvClient.DeleteKey(ctx, path.Base(*kItem.Kid))
			if err != nil {
				return err
			}
		}

		kItems, err = akvClient.GetKeys(ctx, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

// Please keep this function to clean the keys
func cleanAWSKeys(ctx context.Context, awsClient aws.KmsClient) error {
	kItems, err := awsClient.ListKeys(ctx, 0, "")
	if err != nil {
		return err
	}

	for *kItems.Truncated {
		fmt.Println(len(kItems.Keys), *kItems.NextMarker)
		for _, kItem := range kItems.Keys {
			_, err = awsClient.DeleteKey(ctx, *kItem.KeyId)
			if err != nil {
				continue
			}
		}

		kItems, err = awsClient.ListKeys(ctx, 0, *kItems.NextMarker)
		if err != nil {
			return err
		}
	}

	return nil
}
