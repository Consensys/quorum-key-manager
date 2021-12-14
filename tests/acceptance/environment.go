package acceptancetests

import (
	"context"
	"fmt"
	models2 "github.com/consensys/quorum-key-manager/src/aliases/database/models"
	"os"
	"strconv"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	postgresclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker"
	dconfig "github.com/consensys/quorum-key-manager/tests/acceptance/docker/config"
	"github.com/consensys/quorum-key-manager/tests/acceptance/utils"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	hashicorpContainerID = "hashicorp-vault"
	postgresContainerID  = "postgres"
	networkName          = "key-manager"
	localhostPath        = "http://localhost"
	MaxRetries           = 4
	WaitContainerTime    = 15 * time.Second
)

type IntegrationEnvironment struct {
	logger           log.Logger
	hashicorpAddress string
	hashicorpToken   string
	dockerClient     *docker.Client
	postgresClient   *postgresclient.PostgresClient
	baseURL          string
}

func StartEnvironment(ctx context.Context, env *IntegrationEnvironment) error {
	ctx, cancel := context.WithCancel(ctx)

	sig := common.NewSignalListener(func(signal os.Signal) {
		env.logger.Error("interrupt signal has been sent")
		cancel()
	})
	defer sig.Close()

	err := env.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewIntegrationEnvironment() (*IntegrationEnvironment, error) {
	logger, err := zap.NewLogger(zap.NewConfig(zap.DebugLevel, zap.JSONFormat)) // We log panic as we do not need logs
	if err != nil {
		return nil, err
	}

	// Hashicorp
	hashicorpContainer, err := utils.HashicorpContainer(logger)
	if err != nil {
		return nil, err
	}

	// Postgres
	postgresPort := strconv.Itoa(10000 + rand.Intn(10000))
	postgresContainer := postgres.NewDefault().SetPort(postgresPort)

	// Initialize environment container setup
	composition := &dconfig.Composition{
		Containers: map[string]*dconfig.Container{
			hashicorpContainerID: {HashicorpVault: hashicorpContainer},
			postgresContainerID:  {Postgres: postgresContainer},
		},
	}

	dockerClient, err := docker.NewClient(composition, logger)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	envHTTPPort := rand.IntnRange(20000, 28080)
	httpConfig := server.NewDefaultConfig()
	httpConfig.Port = uint32(envHTTPPort)
	postgresCfg := &postgresclient.Config{
		Host:     "127.0.0.1",
		Port:     postgresContainer.Port,
		User:     "postgres",
		Password: postgresContainer.Password,
		Database: "postgres",
	}

	postgresClient, err := postgresclient.New(postgresCfg)
	if err != nil {
		logger.WithError(err).Error("cannot initialize Postgres client")
		return nil, err
	}

	return &IntegrationEnvironment{
		logger:           logger,
		hashicorpAddress: fmt.Sprintf("http://%s:%s", hashicorpContainer.Host, hashicorpContainer.Port),
		hashicorpToken:   hashicorpContainer.RootToken,
		dockerClient:     dockerClient,
		postgresClient:   postgresClient,
		baseURL:          fmt.Sprintf("%s:%d", localhostPath, envHTTPPort),
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

	err = env.dockerClient.WaitTillIsReady(ctx, hashicorpContainerID, WaitContainerTime)
	if err != nil {
		env.logger.WithError(err).Error("could not start vault")
		return err
	}

	// Start Postgres
	err = env.dockerClient.Up(ctx, postgresContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up postgres")
		return err
	}

	err = env.dockerClient.WaitTillIsReady(ctx, postgresContainerID, WaitContainerTime)
	if err != nil {
		env.logger.WithError(err).Error("could not start postgres")
		return err
	}

	err = env.createTables()
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	err := env.dockerClient.Down(ctx, hashicorpContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down vault")
	}

	err = env.dockerClient.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Error("could not down postgres")
	}

	err = env.dockerClient.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not remove network")
	}
}

func (env *IntegrationEnvironment) createTables() error {
	pgCfg, err := env.postgresClient.Config().ToPGOptions()
	if err != nil {
		return err
	}

	db := pg.Connect(pgCfg)
	defer db.Close()

	opts := &orm.CreateTableOptions{
		FKConstraints: true,
	}
	// we create tables for each model
	for _, v := range []interface{}{
		&models.Secret{},
		&models.Key{},
		&models.ETHAccount{},
		&models2.Registry{},
		&models2.Alias{},
	} {
		err = db.Model(v).CreateTable(opts)
		if err != nil {
			return err
		}
	}

	env.logger.Info("tables created successfully from models")
	return nil
}
