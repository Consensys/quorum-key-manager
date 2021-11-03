package src

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	aliasapp "github.com/consensys/quorum-key-manager/src/aliases/app"
	authapp "github.com/consensys/quorum-key-manager/src/auth/app"
	apikey "github.com/consensys/quorum-key-manager/src/infra/api-key/filesystem"
	jwt "github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	tls "github.com/consensys/quorum-key-manager/src/infra/tls/filesystem"
	nodesapp "github.com/consensys/quorum-key-manager/src/nodes/app"
	storesapp "github.com/consensys/quorum-key-manager/src/stores/app"
)

type Config struct {
	HTTP     *server.Config
	Logger   *zap.Config
	Postgres *client.Config
	OIDC     *jwt.Config
	APIKey   *apikey.Config
	TLS      *tls.Config
	Manifest *manifests.Config
}

func New(cfg *Config, logger log.Logger) (*app.App, error) {
	// Create app
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))

	// Register Service Configuration
	err := a.RegisterServiceConfig(&authapp.Config{Manifest: cfg.Manifest, OIDC: cfg.OIDC, TLS: cfg.TLS, APIKey: cfg.APIKey})
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&storesapp.Config{Postgres: cfg.Postgres, Manifest: cfg.Manifest})
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&nodesapp.Config{Manifest: cfg.Manifest})
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&aliasapp.Config{Postgres: cfg.Postgres})
	if err != nil {
		return nil, err
	}

	// Register Services
	err = authapp.RegisterService(a, logger.WithComponent("auth"))
	if err != nil {
		return nil, err
	}

	err = storesapp.RegisterService(a, logger.WithComponent("stores"))
	if err != nil {
		return nil, err
	}

	err = aliasapp.RegisterService(a, logger.WithComponent("aliases"))
	if err != nil {
		return nil, err
	}

	err = nodesapp.RegisterService(a, logger.WithComponent("nodes"))
	if err != nil {
		return nil, err
	}

	return a, nil
}
