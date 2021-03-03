package hashicorp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/mount"

	log "github.com/sirupsen/logrus"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultHashicorpVaultImage = "library/vault:1.6.2"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"

type Vault struct{}

type Config struct {
	Image                 string
	Port                  string
	RootToken             string
	Host                  string
	PluginSourceDirectory string
}

func NewDefault() *Config {
	return &Config{
		Image:     defaultHashicorpVaultImage,
		Port:      defaultHostPort,
		RootToken: defaultRootToken,
		Host:      defaultHost,
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

func (cfg *Config) SetPluginSourceDirectory(dir string) *Config {
	cfg.PluginSourceDirectory = dir
	return cfg
}

func (cfg *Config) DownloadPlugin(filename, version string) (string, error) {
	url := fmt.Sprintf("https://github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/releases/download/%s/orchestrate-hashicorp-vault-plugin", version)
	
	pluginPath := fmt.Sprintf("%s/%s", cfg.PluginSourceDirectory, filename)
	err := downloadPlugin(pluginPath, url)
	if err != nil {
		return "", err
	}
	return pluginPath, nil
}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%v", cfg.RootToken),
		},
		ExposedPorts: nat.PortSet{
			"8200/tcp": struct{}{},
		},
		Tty: true,
		Cmd: []string{"server", "-dev", "-dev-plugin-dir=/vault/plugins", "-log-level=trace"},
	}

	hostConfig := &dockercontainer.HostConfig{
		CapAdd: []string{"IPC_LOCK"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.PluginSourceDirectory,
				Target: "/vault/plugins",
			},
		},
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

	httpClient := http.Client{}

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			resp, err := httpClient.Get(fmt.Sprintf("http://%v:%v/v1/sys/health", cfg.Host, cfg.Port))

			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for Hashicorp Vault service to start")
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for Hashicorp Vault service to be ready")
			default:
				log.WithContext(rctx).Infof("Hashicorp Vault container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}

func downloadPlugin(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = os.Chmod(filepath, 0777)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
