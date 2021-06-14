package localstack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config/localstack"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type Vault struct{}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*localstack.Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	log.Printf("Configuration for localstack %v", cfg)

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
			log.Printf("LocalStack Health on %s", fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port))
			resp, err := httpClient.Get(fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port))

			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for localstack service to start")
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
					log.WithContext(rctx).Infof("localstack container service is ready, running")
					break waitForServiceLoop
				} else {
					log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for localstack service to be ready, status : %s", localStackStatus.Status)
				}
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for localstack service to be ready")
			default:
				log.WithContext(rctx).Infof("localstack container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
