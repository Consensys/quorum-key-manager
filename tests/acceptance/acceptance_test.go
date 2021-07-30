// +build acceptance

package acceptancetests

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/postgres"
	akvkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
	awskey "github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/local"
	akvsecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	awssecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/aws"
	hashicorpsecret "github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
	"github.com/stretchr/testify/suite"
	"os"
	"path"
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
	testSuite.connector = keys.NewConnector(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), db, logger)
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Keys-AKV")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.connector = keys.NewConnector(akvkey.New(s.env.akvClient, logger), db, logger)
	suite.Run(s.T(), testSuite)

	// AWS
	logger = s.env.logger.WithComponent("Keys-AWS")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.connector = keys.NewConnector(awskey.New(s.env.awsKmsClient, db.Keys(), logger), db, logger)
	suite.Run(s.T(), testSuite)

	// Local
	logger = s.env.logger.WithComponent("Keys-Local")
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	testSuite.connector = keys.NewConnector(local.New(hashicorpSecretStore, logger), db, logger)
	suite.Run(s.T(), testSuite)
}

/*func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Eth1-DB"), s.env.postgresClient)

	// Hashicorp
	logger := s.env.logger.WithComponent("Eth1-Hashicorp")
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.connector = eth1.New(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AKV
	logger = s.env.logger.WithComponent("Eth1-AKV")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.connector = eth1.New(akvkey.New(s.env.akvClient, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// AWS
	logger = s.env.logger.WithComponent("Eth1-AWS")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.connector = eth1.New(awskey.New(s.env.awsKmsClient, db.Keys(), logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)

	// Local
	logger = s.env.logger.WithComponent("Eth1-Local")
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, HashicorpSecretMountPoint, logger)
	testSuite.connector = eth1.New(local.New(hashicorpSecretStore, logger), db, logger)
	testSuite.db = db
	suite.Run(s.T(), testSuite)
}*/

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
