// +build acceptance

package acceptancetests

import (
	"context"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	aliasint "github.com/consensys/quorum-key-manager/src/aliases/service/aliases"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	eth "github.com/consensys/quorum-key-manager/src/stores/connectors/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/local"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
	utils2 "github.com/consensys/quorum-key-manager/src/utils/service/utils"
	"github.com/consensys/quorum-key-manager/tests/acceptance/utils"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type acceptanceTestSuite struct {
	suite.Suite
	env                  *IntegrationEnvironment
	auth                 *authorizator.Authorizator
	utils                *utils2.Utilities
	hashicorpKvv2Client  *client.HashicorpVaultClient
	hasicorpPluginClient *client.HashicorpVaultClient
	db                   database.Database
	err                  error
}

func (s *acceptanceTestSuite) SetupSuite() {
	err := StartEnvironment(context.Background(), s.env)
	require.NoError(s.T(), err)

	s.hashicorpKvv2Client, err = client.NewClient(client.NewConfig(&entities.HashicorpConfig{
		MountPoint: "secret",
		Address:    s.env.hashicorpAddress,
	}))
	require.NoError(s.T(), err)

	s.hasicorpPluginClient, err = client.NewClient(client.NewConfig(&entities.HashicorpConfig{
		MountPoint: "quorum",
		Address:    s.env.hashicorpAddress,
	}))
	require.NoError(s.T(), err)

	s.hashicorpKvv2Client.SetToken(s.env.hashicorpToken)
	s.hasicorpPluginClient.SetToken(s.env.hashicorpToken)

	err = s.hasicorpPluginClient.Mount("quorum", &api.MountInput{
		Type:        "plugin",
		Description: "Quorum Hashicorp Vault Plugin",
		Config: api.MountConfigInput{
			ForceNoCache:              true,
			PassthroughRequestHeaders: []string{"X-Vault-Namespace"},
		},
		PluginName: utils.HashicorpPluginFilename,
	})
	require.NoError(s.T(), err)

	s.auth = authorizator.New(authtypes.ListPermissions(), "", s.env.logger)
	s.utils = utils2.New(s.env.logger)
	s.db = postgres.New(s.env.logger, s.env.postgresClient)

	s.env.logger.Info("setup test suite has completed")
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerStore(t *testing.T) {
	env, err := NewIntegrationEnvironment()
	require.NoError(t, err)

	s := new(acceptanceTestSuite)
	s.env = env

	suite.Run(t, s)
}

func (s *acceptanceTestSuite) TestKeyManagerStore_Secrets() {
	storeName := "acceptance_secret_store"
	logger := s.env.logger.WithComponent(storeName)
	db := s.db.Secrets(storeName)
	secretStore := hashicorp.New(s.hashicorpKvv2Client, db, s.env.logger)

	testSuite := new(secretsTestSuite)
	testSuite.env = s.env
	testSuite.db = db
	testSuite.store = secrets.NewConnector(secretStore, db, s.auth, logger)

	suite.Run(s.T(), testSuite)
}

func (s *acceptanceTestSuite) TestKeyManager_Keys() {
	// Hashicorp
	storeName := "acceptance_key_store_hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	db := s.db.Keys(storeName)

	testSuite := new(keysTestSuite)
	testSuite.env = s.env
	testSuite.db = db
	testSuite.store = keys.NewConnector(hashicorpkey.New(s.hasicorpPluginClient, logger), db, s.auth, logger)
	testSuite.utils = s.utils

	suite.Run(s.T(), testSuite)

	// Local
	storeName = "acceptance_key_store_local"
	logger = s.env.logger.WithComponent(storeName)
	db = s.db.Keys(storeName)
	secretsDB := s.db.Secrets(storeName)

	testSuite = new(keysTestSuite)
	testSuite.env = s.env
	testSuite.db = db
	secretStore := hashicorp.New(s.hashicorpKvv2Client, secretsDB, s.env.logger)
	testSuite.utils = s.utils
	testSuite.store = keys.NewConnector(local.New(secretStore, secretsDB, logger), db, s.auth, logger)

	suite.Run(s.T(), testSuite)
}

func (s *acceptanceTestSuite) TestKeyManagerStore_Eth() {
	// Hashicorp
	storeName := "acceptance_ethereum_store_hashicorp"
	logger := s.env.logger.WithComponent(storeName)
	db := s.db.ETHAccounts(storeName)

	testSuite := new(ethTestSuite)
	testSuite.env = s.env
	testSuite.db = db
	testSuite.store = eth.NewConnector(hashicorpkey.New(s.hasicorpPluginClient, logger), db, s.auth, logger)
	testSuite.utils = s.utils

	suite.Run(s.T(), testSuite)

	// Local
	storeName = "acceptance_ethereum_store_local"
	logger = s.env.logger.WithComponent(storeName)
	db = s.db.ETHAccounts(storeName)
	secretsDB := s.db.Secrets(storeName)

	testSuite = new(ethTestSuite)
	testSuite.env = s.env
	testSuite.db = db
	testSuite.utils = s.utils
	testSuite.store = eth.NewConnector(local.New(hashicorp.New(s.hashicorpKvv2Client, secretsDB, logger), secretsDB, logger), db, s.auth, logger)

	suite.Run(s.T(), testSuite)
}

func (s *acceptanceTestSuite) TestKeyManagerAliases() {
	db := aliaspg.NewDatabase(s.env.postgresClient, s.env.logger).Alias()

	testSuite := new(aliasStoreTestSuite)
	testSuite.env = s.env
	testSuite.srv = aliasint.New(db, s.env.logger)
	testSuite.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	suite.Run(s.T(), testSuite)
}
