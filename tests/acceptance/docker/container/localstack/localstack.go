package localstack

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"net/http"
	"strings"
	"time"

	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config/localstack"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type Vault struct {
	logger log.Logger
}

func New(logger log.Logger) *Vault {
	return &Vault{
		logger: logger,
	}
}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*localstack.Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("AWS_DEFAULT_REGION=%v", cfg.Region),
			fmt.Sprintf("SERVICES=%v", strings.Join(cfg.Services, ",")),
		},
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%s/tcp", localstack.DefaultHostPort)): struct{}{},
		},
		Tty:        true,
		Entrypoint: []string{"docker-entrypoint.sh"},
	}

	hostConfig := &dockercontainer.HostConfig{}

	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"4566/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (vault *Vault) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*localstack.Config)
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
			resp, err := httpClient.Get(fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port))

			switch {
			case err != nil:
				vault.logger.WithError(err).Debug("waiting for localstack service to start")
			case resp.StatusCode == http.StatusNotFound:
				//Ready status is encoded in resp Body
				defer resp.Body.Close()
				type LocalStackStatus struct {
					Status string `json:"status"`
				}
				//Found no better way to ensure readiness
				localStackStatus := LocalStackStatus{}
				_ = json.NewDecoder(resp.Body).Decode(&localStackStatus)
				if localStackStatus.Status == "running" {
					vault.logger.Info("localstack container service is ready")
					break waitForServiceLoop
				} else {
					vault.logger.Debug("waiting for localstack service to be ready", "status_code", resp.StatusCode, "status", localStackStatus)
				}
			case resp.StatusCode != http.StatusOK:
				vault.logger.Debug("waiting for localstack service to be ready", "status_code", resp.StatusCode)
			default:
				vault.logger.Info("localstack container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
