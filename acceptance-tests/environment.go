package integrationtests

import (
	"context"
	"fmt"
	akvclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/hashicorp/vault/api"
)

const hashicorpContainerID = "hashicorp-vault"
const networkName = "key-manager"

type IntegrationEnvironment struct {
	ctx             context.Context
	logger          *log.Logger
	hashicorpClient *hashicorpclient.HashicorpVaultClient
	akvClient       *akvclient.AzureClient
	dockerClient    *docker.Client
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

	hashicorpContainer, err := HashicorpContainer(ctx)
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

	hashicorpAddr := fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)
	hashicorpClient, err := hashicorpclient.NewClient(hashicorpclient.NewBaseConfig(hashicorpAddr, hashicorpContainer.RootToken))
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}

	akvClient, err := akvclient.NewClient(akvclient.NewConfig(os.Getenv("AKV_VAULT_NAME"), os.Getenv("AKV_TENANT_ID"), os.Getenv("AKV_CLIENT_ID"), os.Getenv("AKV_CLIENT_SECRET")))
	if err != nil {
		logger.WithError(err).Error("cannot initialize akv client")
		return nil, err
	}

	// TODO Add key-manager init

	return &IntegrationEnvironment{
		ctx:             ctx,
		logger:          logger,
		hashicorpClient: hashicorpClient,
		akvClient:       akvClient,
		dockerClient:    dockerClient,
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
		PluginName: hashicorpPluginFilename,
	})
	if err != nil {
		env.logger.WithError(err).Error("failed to mount (enable) orchestrate vault plugin")
		return err
	}

	// TODO Add key-manager START and WAIT_FOR_READY

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	// TODO Add key-manager STOP

	err := env.dockerClient.Down(ctx, hashicorpContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down vault")
	}

	err = env.dockerClient.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not remove network")
	}
}
