package app

import (
	"context"
	"crypto/x509"

	"github.com/consensys/quorum-key-manager/src/auth/api/middlewares"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/manager"
	"github.com/consensys/quorum-key-manager/src/auth/service/authenticator"
	apikey "github.com/consensys/quorum-key-manager/src/infra/api-key/filesystem"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	tls "github.com/consensys/quorum-key-manager/src/infra/tls/filesystem"
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
	ctx := context.Background()

	var authMiddleware alice.Constructor
	authEnabled := cfg.OIDC != nil || cfg.APIKey != nil || cfg.TLS != nil
	if authEnabled {
		var jwtValidator jwt.Validator
		var apikeyClaims map[string]*entities.UserClaims
		var rootCAs *x509.CertPool
		var err error

		if cfg.OIDC != nil {
			jwtValidator, err = jose.New(cfg.OIDC)
			if err != nil {
				return nil, err
			}

			logger.Info("JWT authentication enabled")
		}

		if cfg.APIKey != nil {
			apiKeyReader, err := apikey.New(cfg.APIKey)
			if err != nil {
				return nil, err
			}

			apikeyClaims, err = apiKeyReader.Load(ctx)
			if err != nil {
				return nil, err
			}

			logger.Info("API key authentication enabled")
		}

		if cfg.TLS != nil {
			tlsReader, err := tls.New(cfg.TLS)
			if err != nil {
				return nil, err
			}

			rootCAs, err = tlsReader.Load(ctx)
			if err != nil {
				return nil, err
			}

			logger.Info("TLS authentication enabled")
		}

		authMiddleware = middlewares.NewAuth(authenticator.New(jwtValidator, apikeyClaims, rootCAs, logger)).Middleware
	} else {
		authMiddleware = middlewares.WildcardMiddleware
	}

	composition := alice.New(
		middlewares.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		authMiddleware,
	)

	return &composition, nil
}
