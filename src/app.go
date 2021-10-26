package src

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	aliasapp "github.com/consensys/quorum-key-manager/src/aliases/app"
	"github.com/consensys/quorum-key-manager/src/auth"
	app3 "github.com/consensys/quorum-key-manager/src/auth/app"
	"github.com/consensys/quorum-key-manager/src/infra/http/accesslog"
	"github.com/consensys/quorum-key-manager/src/infra/http/middlewares/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	app2 "github.com/consensys/quorum-key-manager/src/nodes/app"
	stores "github.com/consensys/quorum-key-manager/src/stores/app"
	"github.com/justinas/alice"
)

type Config struct {
	HTTP     *server.Config
	Logger   *log.Config
	Postgres *client.Config
	Auth     *auth.Config
	Manifest *manifests.Config
}

func New(cfg *Config, logger log.Logger) (*app.App, error) {
	// Create app
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))

	// Register Service Configuration
	err := a.RegisterServiceConfig(cfg.Auth)
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&stores.Config{Postgres: cfg.Postgres, Manifest: cfg.Manifest})
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&app2.Config{Manifest: cfg.Manifest})
	if err != nil {
		return nil, err
	}

	// Register Services
	err = app3.RegisterService(a, logger.WithComponent("auth"))
	if err != nil {
		return nil, err
	}

	err = stores.RegisterService(a, logger.WithComponent("stores"))
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&aliasapp.Config{Postgres: cfg.Postgres})
	if err != nil {
		return nil, err
	}

	err = aliasapp.RegisterService(a, logger.WithComponent("aliases"))
	if err != nil {
		return nil, err
	}

	err = app2.RegisterService(a, logger.WithComponent("nodes"))
	if err != nil {
		return nil, err
	}

	// Set Middleware
	compositionMiddleware, err := createMiddlewares(cfg, logger.WithComponent("accesslog"))
	if err != nil {
		return nil, err
	}

	err = a.SetMiddleware(compositionMiddleware.Then)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func createMiddlewares(cfg *Config, logger log.Logger) (*alice.Chain, error) {
	compositionMiddleware := alice.Chain{}

	// access log middleware
	// TODO: Make accesslog middleware configurable (at least enable/disable)
	compositionMiddleware = compositionMiddleware.Append(accesslog.NewMiddleware(logger).Handler)

	if cfg.Auth.OIDC.IssuerURL != "" {
		jwtMiddleware, err := jwt.New(cfg.Auth.OIDC)
		if err != nil {
			errMessage := "failed to create jwt middleware"
			logger.WithError(err).Error(errMessage, "issuer_url", cfg.Auth.OIDC.IssuerURL)
			return nil, errors.ConfigError(errMessage)
		}

		compositionMiddleware = compositionMiddleware.Append(jwtMiddleware.Handler)
	}

	return &compositionMiddleware, nil
}
