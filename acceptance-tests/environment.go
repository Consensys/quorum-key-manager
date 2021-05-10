package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
	"gopkg.in/yaml.v2"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/utils"
	keymanager "github.com/ConsenSysQuorum/quorum-key-manager/src"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/hashicorp"
	akvclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/http"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/akv"
	akv2 "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/hashicorp/vault/api"
)

const (
	hashicorpContainerID      = "hashicorp-vault"
	networkName               = "key-manager"
	localhostPath             = "http://localhost"
	HashicorpSecretStoreName  = "HashicorpSecrets"
	HashicorpKeyStoreName     = "HashicorpKeys"
	HashicorpSecretMountPoint = "secret"
	HashicorpKeyMountPoint    = "orchestrate"
	AKVSecretStoreName        = "AKVSecrets"
	AKVKeyStoreName           = "AKVKeys"
	AKVSpecENV                = "AKV_ENVIRONMENT"
)

type IntegrationEnvironment struct {
	ctx             context.Context
	logger          *log.Logger
	hashicorpClient hashicorp2.VaultClient
	akvClient       akv2.Client
	dockerClient    *docker.Client
	keyManager      *keymanager.App
	baseURL         string
	Cancel          context.CancelFunc
	tmpYml          string
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
	logger := log.NewDefaultLogger().WithContext(ctx)

	hashicorpContainer, err := utils.HashicorpContainer(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize environment container setup
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			hashicorpContainerID: {
				HashicorpVault: hashicorpContainer,
			},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	akvSpecStr := os.Getenv(AKVSpecENV)
	akvKeySpecs := akv.KeySpecs{}
	akvSecretSpecs := akv.KeySpecs{}
	_ = json.Unmarshal([]byte(akvSpecStr), &akvKeySpecs)
	_ = json.Unmarshal([]byte(akvSpecStr), &akvSecretSpecs)

	envHTTPPort := rand.IntnRange(20000, 28080)
	hashicorpAddr := fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)
	tmpYml, err := newTemporalManifestYml(&manifest.Manifest{
		Kind: types.HashicorpSecrets,
		Name: HashicorpSecretStoreName,
		Specs: &hashicorp.SecretSpecs{
			Token:      hashicorpContainer.RootToken,
			MountPoint: HashicorpSecretMountPoint,
			Address:    hashicorpAddr,
			Namespace:  "",
		},
	}, &manifest.Manifest{
		Kind: types.HashicorpKeys,
		Name: HashicorpKeyStoreName,
		Specs: &hashicorp.KeySpecs{
			MountPoint: HashicorpKeyMountPoint,
			Address:    hashicorpAddr,
			Token:      hashicorpContainer.RootToken,
			Namespace:  "",
		},
	}, &manifest.Manifest{
		Kind:  types.AKVSecrets,
		Name:  AKVSecretStoreName,
		Specs: akvSecretSpecs,
	}, &manifest.Manifest{
		Kind:  types.AKVKeys,
		Name:  AKVKeyStoreName,
		Specs: akvKeySpecs,
	})

	if err != nil {
		logger.WithError(err).Error("cannot create keymanager manifest")
		return nil, err
	} else {
		logger.WithField("path", tmpYml).Info("new temporal manifest created")
	}

	httpConfig := http.NewDefaultConfig()
	httpConfig.Port = uint32(envHTTPPort)
	keyManager, err := newKeyManager(&keymanager.Config{
		HTTP:         httpConfig,
		ManifestPath: tmpYml,
	}, logger)
	if err != nil {
		logger.WithError(err).Error("cannot initialize keymanager server")
		return nil, err
	}

	// Hashicorp client for direct integration tests
	hashicorpClient, err := hashicorpclient.NewClient(hashicorpclient.NewBaseConfig(hashicorpAddr, hashicorpContainer.RootToken, ""))
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}

	akvClient, err := akvclient.NewClient(akvclient.NewConfig(
		akvKeySpecs.VaultName,
		akvKeySpecs.TenantID,
		akvKeySpecs.ClientID,
		akvKeySpecs.ClientSecret,
	))
	if err != nil {
		logger.WithError(err).Error("cannot initialize akv client")
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &IntegrationEnvironment{
		ctx:             ctx,
		logger:          logger,
		hashicorpClient: hashicorpClient,
		akvClient:       akvClient,
		dockerClient:    dockerClient,
		keyManager:      keyManager,
		baseURL:         fmt.Sprintf("%s:%d", localhostPath, envHTTPPort),
		Cancel:          cancel,
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

	go func() {
		err = env.keyManager.Start(ctx)
		if err != nil {
			env.logger.WithError(err).Error("failed to start key manager")
			env.Cancel()
		}
	}()

	// TODO: Implement WaitFor functions based on ready endpoint
	time.Sleep(2 * time.Second)

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

	err = os.Remove(env.tmpYml)
	if err != nil {
		env.logger.WithError(err).Error("cannot remove temporal yml file")
	}
}

func newTemporalManifestYml(manifests ...*manifest.Manifest) (string, error) {
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

func newKeyManager(cfg *keymanager.Config, logger *log.Logger) (*keymanager.App, error) {
	return keymanager.New(cfg, logger)
}
