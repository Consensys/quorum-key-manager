package integrationtests

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/common/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/infra/vault/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker"
	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker/config"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/pflag"
)

const hashicorpContainerID = "hashicorp-vault"
const networkName = "key-manager"
const vaultTokenFilePrefix = "orchestrate_vault_token_"

type IntegrationEnvironment struct {
	ctx             context.Context
	logger          *log.Logger
	hashicorpClient *hashicorp.HashicorpVaultClient
	dockerClient    *docker.Client
}

type TestSuiteEnv interface {
	Start(ctx context.Context) error
}

func StartEnvironment(ctx context.Context, env TestSuiteEnv) (gerr error) {
	ctx, cancel := context.WithCancel(ctx)

	sig := NewSignalListener(func(signal os.Signal) {
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
	logger := log.NewLogger().WithContext(ctx)

	hashicorpContainer, err := HashicorpContainer(ctx)
	if err != nil {
		return nil, err
	}

	hashicorpTokenFileName, err := generateHashicorpTokenFile(hashicorpContainer.RootToken)
	if err != nil {
		logger.WithError(err).Error("cannot generate vault token file")
		return nil, err
	}

	// Initialize environment flags
	flgs := pflag.NewFlagSet("integration-tests", pflag.ContinueOnError)

	hashicorp.Flags(flgs)
	log.Flags(flgs)
	args := []string{
		fmt.Sprintf("--%s=%s", hashicorp.HashicorpAddrFlag,
			fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)),
		fmt.Sprintf("--%s=%s", hashicorp.HashicorpTokenFilePathFlag, hashicorpTokenFileName),
		"--log-level=panic",
	}

	err = flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
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

	hashicorpCfg := hashicorp.ConfigFromViper()
	hashicorpCfg.Renewal = false
	hashicorpClient, err := hashicorp.NewVaultClient(ctx, hashicorpCfg)
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}
	hashicorpClient.Client().SetToken(hashicorpContainer.RootToken)

	// TODO Add key-manager init

	return &IntegrationEnvironment{
		ctx:             ctx,
		logger:          logger,
		hashicorpClient: hashicorpClient,
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
