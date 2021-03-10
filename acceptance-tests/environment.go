package integrationtests

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/cmd/flags"
	app "github.com/ConsenSysQuorum/quorum-key-manager/src"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/utils"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/hashicorp/vault/api"
)

const hashicorpContainerID = "hashicorp-vault"
const networkName = "key-manager"
const localhostPath = "http://localhost:"

type IntegrationEnvironment struct {
	ctx             context.Context
	logger          *log.Logger
	hashicorpClient hashicorp.VaultClient
	dockerClient    *docker.Client
	baseURL         string
	keyManager      *app.App
}

type TestSuiteEnv interface {
	Start(ctx context.Context) error
}

func StartEnvironment(ctx context.Context, env TestSuiteEnv) (gerr error) {
	ctx, cancel := context.WithCancel(ctx)

	sig := utils.NewSignalListener(func(signal os.Signal) {
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
	hashicorpClient, err := client.NewClient(client.NewBaseConfig(hashicorpAddr, "acceptance-test", hashicorpContainer.RootToken))
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}

	flgs := pflag.NewFlagSet("key-manager-acceptance-test", pflag.ContinueOnError)
	flags.HashicorpFlags(flgs)
	flags.LoggerFlags(flgs)
	args := []string{
		"--hashicorp-addr=http://" + hashicorpContainer.Host + ":" + hashicorpContainer.Port,
		"--hashicorp-token=" + hashicorpContainer.RootToken,
		"--log-level=debug",
	}

	err = flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
		return nil, err
	}

	return &IntegrationEnvironment{
		ctx:             ctx,
		logger:          logger,
		hashicorpClient: hashicorpClient,
		dockerClient:    dockerClient,
		baseURL:         localhostPath + "8080",
		keyManager:      app.New(flags.NewAppConfig(viper.GetViper())),
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

	err = env.keyManager.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("failed to start key manager")
		return err
	}

	// TODO: Wait for service to be ready instead of sleeping
	time.Sleep(5 * time.Second)

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
}
