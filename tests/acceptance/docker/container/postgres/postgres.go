package postgres

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	postgresclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	postgresConfig "github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/postgres"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/go-pg/pg/v10"
	"time"
)

type Postgres struct {
	logger log.Logger
}

func New(logger log.Logger) *Postgres {
	return &Postgres{
		logger: logger,
	}
}

func (p *Postgres) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*postgresConfig.Config)
	if !ok {
		errMessage := fmt.Sprintf("invalid configuration type (expected %T but got %T)", cfg, configuration)
		p.logger.Error(errMessage)
		return nil, nil, nil, fmt.Errorf(errMessage)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%v", cfg.Password),
		},
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"5432/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (p *Postgres) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*postgresConfig.Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	posgresClientCfg := &postgresclient.Config{
		Host:     "127.0.0.1",
		Port:     cfg.Port,
		User:     "postgres",
		Password: cfg.Password,
		Database: "postgres",
	}
	pgCfg, _ := posgresClientCfg.ToPGOptions()
	db := pg.Connect(pgCfg)
	defer db.Close()

	retryT := time.NewTicker(time.Second)
	defer retryT.Stop()

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			_, err := db.Exec("SELECT 1")
			if err != nil {
				p.logger.WithError(err).Warn("waiting for PostgreSQL service to start")
			} else {
				p.logger.Info("PostgreSQL container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
