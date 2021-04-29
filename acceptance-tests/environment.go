package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/utils"
	keymanager "github.com/ConsenSysQuorum/quorum-key-manager/src"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
	akvclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/http"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/cmd/flags"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/akv"
	akv2 "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/hashicorp/vault/api"
)

const (
	hashicorpContainerID = "hashicorp-vault"
	networkName          = "key-manager"
	localhostPath        = "http://localhost"
	SecretStoreName      = "HashicorpSecrets"
	KeyStoreName         = "HashicorpKeys"
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

	envHTTPPort := rand.IntnRange(20000, 28080)
	hashicorpAddr := fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port)
	hashicorpSecretSpecs := &hashicorp.SecretSpecs{
		Token:      hashicorpContainer.RootToken,
		MountPoint: "secret",
		Address:    hashicorpAddr,
		Namespace:  "",
	}
	hashicorpKeySpecs := &hashicorp.KeySpecs{
		MountPoint: "orchestrate",
		Address:    hashicorpAddr,
		Token:      hashicorpContainer.RootToken,
		Namespace:  "",
	}
	keyManager := newKeyManager(logger, hashicorpSecretSpecs, hashicorpKeySpecs, uint32(envHTTPPort))

	// Hashicorp client for direct integration tests
	hashicorpClient, err := hashicorpclient.NewClient(hashicorpclient.NewBaseConfig(hashicorpSecretSpecs.Address, hashicorpSecretSpecs.Token, ""))
	if err != nil {
		logger.WithError(err).Error("cannot initialize hashicorp vault client")
		return nil, err
	}

	akvSpecStr := os.Getenv(flags.AKVEnvironmentEnv)
	specs := akv.Specs{}
	_ = json.Unmarshal([]byte(akvSpecStr), &specs)

	akvClient, err := akvclient.NewClient(akvclient.NewConfig(
		specs.VaultName,
		specs.TenantID,
		specs.ClientID,
		specs.ClientSecret,
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
	return nil
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
	return
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

func newKeyManager(
	logger *log.Logger,
	hashicorpSecretStoreSpecs *hashicorp.SecretSpecs,
	hashicorpKeyStoreSpecs *hashicorp.KeySpecs,
	port uint32,
) *keymanager.App {
	hashicorpSecretSpecsRaw, _ := json.Marshal(hashicorpSecretStoreSpecs)
	hashicorpSecretManifest := &manifest.Manifest{
		Kind:    types.HashicorpSecrets,
		Name:    SecretStoreName,
		Version: "0.0.0",
		Specs:   hashicorpSecretSpecsRaw,
	}

	hashicorpKeySpecsRaw, _ := json.Marshal(hashicorpKeyStoreSpecs)
	hashicorpKeyManifest := &manifest.Manifest{
		Kind:    types.HashicorpKeys,
		Name:    KeyStoreName,
		Version: "0.0.0",
		Specs:   hashicorpKeySpecsRaw,
	}

	httpConfig := http.NewDefaultConfig()
	httpConfig.Port = port
	cfg := &keymanager.Config{
		HTTP: httpConfig,
		Manifests: []*manifest.Manifest{
			hashicorpSecretManifest,
			hashicorpKeyManifest,
		},
	}

	return keymanager.New(cfg, logger)
}
