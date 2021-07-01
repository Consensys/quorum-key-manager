package acceptancetests

import (
	"context"
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/pkg/log/zap"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws"
	"io/ioutil"
	"os"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/app"
	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/server"
	keymanager "github.com/consensysquorum/quorum-key-manager/src"
	manifestsmanager "github.com/consensysquorum/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
	akv2 "github.com/consensysquorum/quorum-key-manager/src/stores/infra/akv"
	akvclient "github.com/consensysquorum/quorum-key-manager/src/stores/infra/akv/client"
	awsclient "github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws/client"
	hashicorp2 "github.com/consensysquorum/quorum-key-manager/src/stores/infra/hashicorp"
	hashicorpclient "github.com/consensysquorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/consensysquorum/quorum-key-manager/src/stores/types"
	"github.com/consensysquorum/quorum-key-manager/tests"
	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker"
	dconfig "github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config"
	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/utils"
	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	hashicorpContainerID      = "hashicorp-vault"
	networkName               = "key-manager"
	localhostPath             = "http://localhost"
	HashicorpKeyStoreName     = "HashicorpKeys"
	HashicorpSecretMountPoint = "secret"
	HashicorpKeyMountPoint    = "orchestrate"
	AKVKeyStoreName           = "AKVKeys"
	AWSKeyStoreName           = "AWSKeys"
)

type IntegrationEnvironment struct {
	ctx               context.Context
	logger            log.Logger
	hashicorpClient   hashicorp2.VaultClient
	awsSecretsClient  aws.SecretsManagerClient
	awsKmsClient      aws.KmsClient
	akvClient         akv2.Client
	dockerClient      *docker.Client
	keyManager        *app.App
	baseURL           string
	Cancel            context.CancelFunc
	tmpManifestYaml   string
	tmpHashicorpToken string
	cfg               *tests.Config
}

const MaxRetries = 10

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
	logger, err := zap.NewLogger(log.NewConfig(log.ErrorLevel, true, log.ProductionMode))
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

	// Initialize environment container setup
	composition := &dconfig.Composition{
		Containers: map[string]*dconfig.Container{
			hashicorpContainerID: {
				HashicorpVault: hashicorpContainer,
			},
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
			Specs: &hashicorp.KeySpecs{
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

	logger.Info("new temporary manifest created", "path", tmpYml)

	httpConfig := server.NewDefaultConfig()
	httpConfig.Port = uint32(envHTTPPort)
	keyManager, err := keymanager.New(&keymanager.Config{
		HTTP:      httpConfig,
		Manifests: &manifestsmanager.Config{Path: tmpYml},
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
	hashicorpClient.Client().SetToken(hashicorpContainer.RootToken)

	awsConfig := awsclient.NewConfig(testCfg.AwsClient.Region, testCfg.AwsClient.AccessID, testCfg.AwsClient.SecretKey, false)
	akvClient, err := akvclient.NewClient(akvclient.NewConfig(
		testCfg.AkvClient.VaultName,
		testCfg.AkvClient.TenantID,
		testCfg.AkvClient.ClientID,
		testCfg.AkvClient.ClientSecret,
	))
	if err != nil {
		logger.WithError(err).Error("cannot initialize akv client")
		return nil, err
	}

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

	ctx, cancel := context.WithCancel(ctx)
	return &IntegrationEnvironment{
		ctx:               ctx,
		logger:            logger,
		hashicorpClient:   hashicorpClient,
		akvClient:         akvClient,
		awsSecretsClient:  awsSecretsClient,
		awsKmsClient:      awsKeysClient,
		dockerClient:      dockerClient,
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

	err = env.hashicorpClient.Client().Sys().Mount("orchestrate", &api.MountInput{
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
