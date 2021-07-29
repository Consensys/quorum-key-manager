package acceptancetests

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
	"github.com/consensys/quorum-key-manager/src/infra/akv/client"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	awsclient "github.com/consensys/quorum-key-manager/src/infra/aws/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	hashicorpclient "github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	postgresclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	keymanager "github.com/consensys/quorum-key-manager/src"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	hashicorpmanager "github.com/consensys/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/types"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker"
	dconfig "github.com/consensys/quorum-key-manager/tests/acceptance/docker/config"
	"github.com/consensys/quorum-key-manager/tests/acceptance/utils"
	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	hashicorpContainerID      = "hashicorp-vault"
	postgresContainerID       = "postgres"
	networkName               = "key-manager"
	localhostPath             = "http://localhost"
	HashicorpKeyStoreName     = "HashicorpKeys"
	HashicorpSecretMountPoint = "secret"
	HashicorpKeyMountPoint    = "orchestrate"
	AKVKeyStoreName           = "AKVKeys"
	AWSKeyStoreName           = "AWSKeys"
	MaxRetries                = 10
)

type IntegrationEnvironment struct {
	ctx               context.Context
	logger            log.Logger
	hashicorpClient   hashicorp.VaultClient
	awsSecretsClient  aws.SecretsManagerClient
	awsKmsClient      aws.KmsClient
	akvClient         akv.Client
	dockerClient      *docker.Client
	postgresClient    *postgresclient.PostgresClient
	keyManager        *app.App
	baseURL           string
	Cancel            context.CancelFunc
	tmpManifestYaml   string
	tmpHashicorpToken string
	cfg               *tests.Config
}

type TestSuiteEnv interface {
	Start(ctx context.Context) error
}

func StartEnvironment(ctx context.Context, env TestSuiteEnv) (gerr error) {
	ctx, cancel := context.WithCancel(ctx)

	sig := common.NewSignalListener(func(signal os.Signal) {
		gerr = fmt.Errorf("interrupt signal has been sent")
		cancel()
	})
	defer sig.Close()

	err := env.Start(ctx)
	if err != nil {
		if gerr == nil {
			return err
		}
	}

	return
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger, err := zap.NewLogger(log.NewConfig(log.ErrorLevel, log.JSONFormat))
	if err != nil {
		return nil, err
	}

	hashicorpContainer, err := utils.HashicorpContainer(logger)
	if err != nil {
		return nil, err
	}
	tmpTokenFile, err := newTmpFile(hashicorpContainer.RootToken)
	if err != nil {
		return nil, err
	}

	postgresPort := strconv.Itoa(10000 + rand.Intn(10000))
	postgresContainer := postgres.NewDefault().SetPort(postgresPort)

	// Initialize environment container setup
	composition := &dconfig.Composition{
		Containers: map[string]*dconfig.Container{
			hashicorpContainerID: {HashicorpVault: hashicorpContainer},
			postgresContainerID:  {Postgres: postgresContainer},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition, logger)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	testCfg, err := tests.NewConfig()
	if err != nil {
		logger.WithError(err).Error("could not load config")
		return nil, err
	}

	envHTTPPort := rand.IntnRange(20000, 28080)
	hashicorpAddr := fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)
	tmpYml, err := newTmpManifestYml(
		&manifest.Manifest{
			Kind: types.HashicorpKeys,
			Name: HashicorpKeyStoreName,
			Specs: &hashicorpmanager.KeySpecs{
				MountPoint: HashicorpKeyMountPoint,
				Address:    hashicorpAddr,
				TokenPath:  tmpTokenFile,
				Namespace:  "",
			},
		},
		&manifest.Manifest{
			Kind:  types.AKVKeys,
			Name:  AKVKeyStoreName,
			Specs: testCfg.AkvKeySpecs(),
		},
		&manifest.Manifest{
			Kind:  types.AWSKeys,
			Name:  AWSKeyStoreName,
			Specs: testCfg.AwsKeySpecs(),
		},
	)

	if err != nil {
		logger.WithError(err).Error("cannot create keymanager manifest")
		return nil, err
	}

	httpConfig := server.NewDefaultConfig()
	httpConfig.Port = uint32(envHTTPPort)
	postgresCfg := &postgresclient.Config{
		Host:     "127.0.0.1",
		Port:     postgresContainer.Port,
		User:     "postgres",
		Password: postgresContainer.Password,
		Database: "postgres",
	}
	keyManager, err := keymanager.New(&keymanager.Config{
		HTTP:      httpConfig,
		Manifests: &manifestsmanager.Config{Path: tmpYml},
		Postgres:  postgresCfg,
		Logger:    log.NewConfig(log.DebugLevel, log.TextFormat),
	}, logger)
	if err != nil {
		logger.WithError(err).Error("cannot initialize Key Manager server")
		return nil, err
	}

	// Hashicorp client for direct integration tests
	hashicorpCfg := hashicorpclient.NewConfig(hashicorpAddr, "")
	hashicorpClient, err := hashicorpclient.NewClient(hashicorpCfg)
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}
	hashicorpClient.SetToken(hashicorpContainer.RootToken)

	akvClient, err := client.NewClient(client.NewConfig(
		testCfg.AkvClient.VaultName,
		testCfg.AkvClient.TenantID,
		testCfg.AkvClient.ClientID,
		testCfg.AkvClient.ClientSecret,
	))
	if err != nil {
		logger.WithError(err).Error("cannot initialize akv client")
		return nil, err
	}

	awsConfig := awsclient.NewConfig(testCfg.AwsClient.Region, testCfg.AwsClient.AccessID, testCfg.AwsClient.SecretKey, false)
	awsSecretsClient, err := awsclient.NewSecretsClient(awsConfig)
	if err != nil {
		logger.WithError(err).Error("cannot initialize AWS Secret client")
		return nil, err
	}
	awsKeysClient, err := awsclient.NewKmsClient(awsConfig)
	if err != nil {
		logger.WithError(err).Error("cannot initialize AWS KMS client")
		return nil, err
	}

	postgresClient, err := postgresclient.NewClient(postgresCfg)
	if err != nil {
		logger.WithError(err).Error("cannot initialize Postgres client")
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &IntegrationEnvironment{
		ctx:               ctx,
		logger:            logger,
		hashicorpClient:   hashicorpClient,
		akvClient:         akvClient,
		awsSecretsClient:  awsSecretsClient,
		awsKmsClient:      awsKeysClient,
		dockerClient:      dockerClient,
		postgresClient:    postgresClient,
		keyManager:        keyManager,
		baseURL:           fmt.Sprintf("%s:%d", localhostPath, envHTTPPort),
		Cancel:            cancel,
		tmpManifestYaml:   tmpYml,
		tmpHashicorpToken: tmpTokenFile,
		cfg:               testCfg,
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.dockerClient.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start Hashicorp Vault
	err = env.dockerClient.Up(ctx, hashicorpContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up vault container")
		return err
	}

	err = env.dockerClient.WaitTillIsReady(ctx, hashicorpContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start vault")
		return err
	}

	err = env.hashicorpClient.HealthCheck()
	if err != nil {
		env.logger.WithError(err).Error("failed to connect to hashicorp plugin")
		return err
	}

	err = env.hashicorpClient.Mount("orchestrate", &api.MountInput{
		Type:        "plugin",
		Description: "Orchestrate Wallets",
		Config: api.MountConfigInput{
			ForceNoCache:              true,
			PassthroughRequestHeaders: []string{"X-Vault-Namespace"},
		},
		PluginName: utils.HashicorpPluginFilename,
	})
	if err != nil {
		env.logger.WithError(err).Error("failed to mount (enable) orchestrate vault plugin")
		return err
	}

	// Start Postgres
	err = env.dockerClient.Up(ctx, postgresContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up postgres")
		return err
	}

	err = env.dockerClient.WaitTillIsReady(ctx, postgresContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start postgres")
		return err
	}

	err = env.createTables()
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	err := env.keyManager.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Error("failed to stop key manager")
	}

	err = env.dockerClient.Down(ctx, hashicorpContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down vault")
	}

	err = env.dockerClient.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down postgres")
	}

	err = env.dockerClient.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not remove network")
	}

	err = os.Remove(env.tmpManifestYaml)
	if err != nil {
		env.logger.WithError(err).Error("cannot remove temporary manifest yml file")
	}

	err = os.Remove(env.tmpHashicorpToken)
	if err != nil {
		env.logger.WithError(err).Error("cannot remove temporary hashicorp token file")
	}
}

func newTmpFile(data interface{}) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "acceptanceTest_")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bData, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	_, err = file.Write(bData)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func newTmpManifestYml(manifests ...*manifest.Manifest) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "acceptanceTest_")
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := yaml.Marshal(manifests)
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (env *IntegrationEnvironment) createTables() error {
	pgCfg, err := env.postgresClient.Config().ToPGOptions()
	if err != nil {
		return err
	}

	db := pg.Connect(pgCfg)
	defer db.Close()

	opts := &orm.CreateTableOptions{
		FKConstraints: true,
	}

	err = db.Model(&models.Key{}).CreateTable(opts)
	if err != nil {
		return err
	}

	err = db.Model(&models.ETH1Account{}).CreateTable(opts)
	if err != nil {
		return err
	}

	env.logger.Info("tables created successgfully from models")
	return nil
}
