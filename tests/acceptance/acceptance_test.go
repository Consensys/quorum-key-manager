// +build acceptance

package acceptancetests

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/common"
	aliasmanager "github.com/consensys/quorum-key-manager/src/aliases/manager"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/store/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/eth1"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
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
		s.err = err
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

	db := postgres.New(s.env.logger.WithComponent("Secrets-DB"), s.env.postgresClient)
	auth := authorizator.New(types.ListPermissions(), "", s.env.logger)

	// Hashicorp
	storeName := "Secrets-Hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.db = db.Secrets(storeName)
	testSuite.store = secrets.NewConnector(hashicorpsecret.New(s.env.hashicorpClient, testSuite.db, HashicorpSecretMountPoint, logger), db.Secrets(storeName), auth, logger)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManager_Keys() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Keys-DB"), s.env.postgresClient)
	auth := authorizator.New(types.ListPermissions(), "", s.env.logger)

	// Hashicorp
	storeName := "Keys-Hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.db = db.Keys(storeName)
	testSuite.store = keys.NewConnector(hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger), db.Keys(storeName), auth, logger)
	suite.Run(s.T(), testSuite)

	// Local
	storeName = "Keys-Local"
	logger = s.env.logger.WithComponent(storeName)
	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.db = db.Keys(storeName)
	secretsDB := db.Secrets(storeName)
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, secretsDB, HashicorpSecretMountPoint, logger)
	testSuite.store = keys.NewConnector(local.New(hashicorpSecretStore, logger), db.Keys(storeName), auth, logger)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerStore_Eth1() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	db := postgres.New(s.env.logger.WithComponent("Eth1-DB"), s.env.postgresClient)
	auth := authorizator.New(types.ListPermissions(), "", s.env.logger)

	// Hashicorp
	storeName := "Eth1-Hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	hashicorpStore := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	testSuite := new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(hashicorpStore, db.ETH1Accounts(storeName), auth, logger)
	testSuite.db = db.ETH1Accounts(storeName)
	suite.Run(s.T(), testSuite)

	// Local
	storeName = "Eth1-Local-Hashicorp"
	logger = s.env.logger.WithComponent(storeName)
	secretsDB := db.Secrets(storeName)
	hashicorpSecretStore := hashicorpsecret.New(s.env.hashicorpClient, secretsDB, HashicorpSecretMountPoint, logger)
	localStore := local.New(hashicorpSecretStore, logger)
	testSuite = new(eth1TestSuite)
	testSuite.env = s.env
	testSuite.store = eth1.NewConnector(localStore, db.ETH1Accounts(storeName), auth, logger)
	testSuite.db = db.ETH1Accounts(storeName)
	suite.Run(s.T(), testSuite)
}

func (s *storeTestSuite) TestKeyManagerAliases() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(aliasStoreTestSuite)
	testSuite.env = s.env
	testSuite.srv = aliasmanager.New(aliaspg.NewDatabase(s.env.postgresClient))
	randSrc := rand.NewSource(time.Now().UnixNano())
	testSuite.rand = rand.New(randSrc)
	suite.Run(s.T(), testSuite)
}
