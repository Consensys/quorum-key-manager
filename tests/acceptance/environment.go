package acceptancetests

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	manifest2 "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/akv/client"
	client2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/aws/client"
	hashicorp3 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp"
	client3 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/hashicorp"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/types"
	"io/ioutil"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	keymanager "github.com/ConsenSysQuorum/quorum-key-manager/src"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/docker"
	dconfig "github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/utils"
	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	hashicorpContainerID      = "hashicorp-vault"
	localStackContainerID     = "localstack"
	networkName               = "key-manager"
	localhostPath             = "http://localhost"
	HashicorpSecretStoreName  = "HashicorpSecrets"
	HashicorpKeyStoreName     = "HashicorpKeys"
	HashicorpSecretMountPoint = "secret"
	HashicorpKeyMountPoint    = "orchestrate"
	AKVSecretStoreName        = "AKVSecrets"
	AKVKeyStoreName           = "AKVKeys"
)

type IntegrationEnvironment struct {
	ctx               context.Context
	logger            *log.Logger
	hashicorpClient   hashicorp3.VaultClient
	awsSecretsClient  *client2.AwsSecretsClient
	awsKmsClient      *client2.AwsKmsClient
	akvClient         akv.Client
	dockerClient      *docker.Client
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
	logger := log.DefaultLogger().WithContext(ctx)

	hashicorpContainer, err := utils.HashicorpContainer(ctx)
	if err != nil {
		return nil, err
	}

	tmpTokenFile, err := newTmpFile(hashicorpContainer.RootToken)
	if err != nil {
		return nil, err
	}

	localstackContainer, err := utils.LocalstackContainer(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize environment container setup
	composition := &dconfig.Composition{
		Containers: map[string]*dconfig.Container{
			hashicorpContainerID: {
				HashicorpVault: hashicorpContainer,
			},
			localStackContainerID: {
				LocalstackVault: localstackContainer,
			},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	testCfg, err := tests.NewConfig()
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	envHTTPPort := rand.IntnRange(20000, 28080)
	hashicorpAddr := fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)
	tmpYml, err := newTmpManifestYml(&manifest2.Manifest{
		Kind: types2.HashicorpSecrets,
		Name: HashicorpSecretStoreName,
		Specs: &hashicorp2.SecretSpecs{
			Token:      hashicorpContainer.RootToken,
			MountPoint: HashicorpSecretMountPoint,
			Address:    hashicorpAddr,
			Namespace:  "",
		},
	}, &manifest2.Manifest{
		Kind: types2.HashicorpKeys,
		Name: HashicorpKeyStoreName,
		Specs: &hashicorp2.KeySpecs{
			MountPoint: HashicorpKeyMountPoint,
			Address:    hashicorpAddr,
			TokenPath:  tmpTokenFile,
			Namespace:  "",
		},
	}, &manifest2.Manifest{
		Kind:  types2.AKVSecrets,
		Name:  AKVSecretStoreName,
		Specs: testCfg.AkvSecretSpecs(),
	}, &manifest2.Manifest{
		Kind:  types2.AKVKeys,
		Name:  AKVKeyStoreName,
		Specs: testCfg.AkvKeySpecs(),
	})

	if err != nil {
		logger.WithError(err).Error("cannot create keymanager manifest")
		return nil, err
	}

	logger.WithField("path", tmpYml).Info("new temporal manifest created")

	httpConfig := server.NewDefaultConfig()
	httpConfig.Port = uint32(envHTTPPort)
	keyManager, err := newKeyManager(&keymanager.Config{
		HTTP:      httpConfig,
		Manifests: &manager.Config{Path: tmpYml},
	}, logger)
	if err != nil {
		logger.WithError(err).Error("cannot initialize Key Manager server")
		return nil, err
	}

	// Hashicorp client for direct integration tests
	hashicorpCfg := client3.NewConfig(hashicorpAddr, "")
	hashicorpClient, err := client3.NewClient(hashicorpCfg)
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}
	hashicorpClient.Client().SetToken(hashicorpContainer.RootToken)

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

	localstackAddr := fmt.Sprintf("http://%s:%s", localstackContainer.Host, localstackContainer.Port)
	awsSecretsClient, err := client2.NewSecretsClient(client2.NewIntegrationConfig("eu-west-3", localstackAddr))
	awsKeysClient, err := client2.NewKmsClient(client2.NewIntegrationConfig("eu-west-3", localstackAddr))
	if err != nil {
		logger.WithError(err).Error("cannot initialize aws client")
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

	// Start localstack container
	err = env.dockerClient.Up(ctx, localStackContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up localstack container")
		return err
	}

	err = env.dockerClient.WaitTillIsReady(ctx, localStackContainerID, 120*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start localstack")
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

	err = env.dockerClient.Down(ctx, localStackContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down localstack")
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

func newTmpManifestYml(manifests ...*manifest2.Manifest) (string, error) {
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

func newKeyManager(cfg *keymanager.Config, logger *log.Logger) (*app.App, error) {
	return keymanager.New(cfg, logger)
}
