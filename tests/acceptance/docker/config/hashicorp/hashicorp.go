package hashicorp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultHashicorpVaultImage = "consensys/quorum-hashicorp-vault-plugin:v1.1.4"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"
const defaultMountPath = "orchestrate"

type Vault struct{}

type Config struct {
	Image     string
	Host      string
	Port      string
	RootToken string
	MonthPath string
}

func NewDefault() *Config {
	return &Config{
		Image:     defaultHashicorpVaultImage,
		Port:      defaultHostPort,
		RootToken: defaultRootToken,
		Host:      defaultHost,
		MonthPath: defaultMountPath,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetRootToken(rootToken string) *Config {
	cfg.RootToken = rootToken
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (cfg *Config) SetMountPath(mountPath string) *Config {
	cfg.MonthPath = mountPath
	return cfg
}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("PLUGIN_MOUNT_PATH=%v", cfg.MonthPath),
			fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%v", cfg.RootToken),
		},
		ExposedPorts: nat.PortSet{
			"8200/tcp": struct{}{},
		},
		Tty: true,
		Entrypoint: []string{"sh", "-c", `
		(sleep 2 ; vault-init-dev.sh)&
		vault server -dev -dev-plugin-dir=/vault/plugins -dev-listen-address="0.0.0.0:8200" -log-level=debug
		`},
	}

	hostConfig := &dockercontainer.HostConfig{
		CapAdd: []string{"IPC_LOCK"},
	}

	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"8200/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (vault *Vault) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(2 * time.Second)
	defer retryT.Stop()

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())
	serverAddr := "http://" + cfg.Host + ":" + cfg.Port

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			resp, err := httpClient.Get(fmt.Sprintf("%s/v1/sys/health", serverAddr))

			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for Hashicorp Vault service to start")
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for Hashicorp Vault service to be ready")
			default:
				log.WithContext(rctx).Info("hashicorp Vault container service is ready")
				break waitForServiceLoop
			}
		}
	}

	if cerr != nil {
		return cerr
	}

	return nil
}
