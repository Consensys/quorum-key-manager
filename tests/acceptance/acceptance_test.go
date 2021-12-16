// +build acceptance

package acceptancetests

import (
	"context"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	"github.com/consensys/quorum-key-manager/src/aliases/service/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/service/registries"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"github.com/consensys/quorum-key-manager/src/auth/service/roles"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/local"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
	utilsservice "github.com/consensys/quorum-key-manager/src/utils/service/utils"
	"github.com/consensys/quorum-key-manager/tests/acceptance/utils"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type acceptanceTestSuite struct {
	suite.Suite
	env                  *IntegrationEnvironment
	auth                 *authorizator.Authorizator
	utils                *utilsservice.Utilities
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
	s.utils = utilsservice.New(s.env.logger)
	s.db = postgres.New(s.env.logger, s.env.postgresClient)

	s.env.logger.Info("setup test suite has completed")
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManager(t *testing.T) {
	env, err := NewIntegrationEnvironment()
	require.NoError(t, err)

	s := new(acceptanceTestSuite)
	s.env = env

	suite.Run(t, s)
}

func (s *acceptanceTestSuite) TestSecrets() {
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

func (s *acceptanceTestSuite) TestKeys() {
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

func (s *acceptanceTestSuite) TestEthereum() {
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

func (s *acceptanceTestSuite) TestAliases() {
	aliasRepository := aliaspg.NewAlias(s.env.postgresClient)
	registryRepository := aliaspg.NewRegistry(s.env.postgresClient)

	rolesService := roles.New(s.env.logger)

	testSuite := new(aliasStoreTestSuite)
	testSuite.env = s.env
	testSuite.aliasService = aliases.New(aliasRepository, registryRepository, rolesService, s.env.logger)
	testSuite.registryService = registries.New(registryRepository, rolesService, s.env.logger)

	suite.Run(s.T(), testSuite)
}
