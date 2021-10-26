package app

import (
	"github.com/consensys/quorum-key-manager/src/auth/api/middlewares"
	"github.com/consensys/quorum-key-manager/src/auth/manager"
	"github.com/consensys/quorum-key-manager/src/auth/service/authenticator"
	"github.com/consensys/quorum-key-manager/src/infra/api-key/csv"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	"github.com/justinas/alice"

	"github.com/consensys/quorum-key-manager/pkg/app"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load configuration
	cfg := new(Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	manifestReader, err := manifestreader.New(cfg.Manifest)
	if err != nil {
		return err
	}

	compositionMiddleware, err := createMiddlewares(cfg, logger)
	if err != nil {
		return err
	}

	err = a.SetMiddleware(compositionMiddleware.Then)
	if err != nil {
		return err
	}

	policyMngr := manager.New(manifestReader, logger)
	err = a.RegisterService(policyMngr)
	if err != nil {
		return err
	}

	return nil
}

func createMiddlewares(cfg *Config, logger log.Logger) (*alice.Chain, error) {
	var authMiddleware alice.Constructor
	authEnabled := cfg.OIDC != nil && cfg.APIKey != nil && cfg.TLS != nil
	if authEnabled {
		jwtValidator, err := jose.New(cfg.OIDC)
		if err != nil {
			return nil, err
		}

		csvReader, err := csv.New(cfg.APIKey)
		if err != nil {
			return nil, err
		}

		authMiddleware = middlewares.NewAuth(authenticator.New(jwtValidator, csvReader, logger)).Middleware
	} else {
		authMiddleware = middlewares.WildcardMiddleware
	}

	composition := alice.New(
		middlewares.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		authMiddleware,
	)

	return &composition, nil
}
