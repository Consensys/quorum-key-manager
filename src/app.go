package src

import (
	"context"
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	aliasapp "github.com/consensys/quorum-key-manager/src/aliases/app"
	authapp "github.com/consensys/quorum-key-manager/src/auth/app"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/api-key/csv"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	tls "github.com/consensys/quorum-key-manager/src/infra/tls/filesystem"
	nodesapp "github.com/consensys/quorum-key-manager/src/nodes/app"
	storesapp "github.com/consensys/quorum-key-manager/src/stores/app"
)

type Config struct {
	HTTP     *server.Config
	Logger   *zap.Config
	Postgres *client.Config
	OIDC     *jose.Config
	APIKey   *csv.Config
	TLS      *tls.Config
	Manifest *manifestreader.Config
}

func New(ctx context.Context, cfg *Config, logger log.Logger) (*app.App, error) {
	// Infra layer
	pgClient, err := client.New(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	jwtValidator, err := getJWTValidator(cfg.OIDC, logger)
	if err != nil {
		return nil, err
	}
	apikeyClaims, err := getAPIKeys(ctx, cfg.APIKey, logger)
	if err != nil {
		return nil, err
	}

	rootCAs, err := getRootCAs(ctx, cfg.TLS, logger)
	if err != nil {
		return nil, err
	}

	manifests, err := getManifests(ctx, cfg.Manifest)
	if err != nil {
		return nil, err
	}

	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))

	// Register Services
	authService, err := authapp.RegisterService(ctx, a, logger.WithComponent("auth"), manifests, jwtValidator, apikeyClaims, rootCAs)
	if err != nil {
		return nil, err
	}

	aliasService, err := aliasapp.RegisterService(a, logger.WithComponent("aliases"), pgClient)
	if err != nil {
		return nil, err
	}

	storesService, err := storesapp.RegisterService(ctx, a, logger.WithComponent("stores"), pgClient, manifests, authService)
	if err != nil {
		return nil, err
	}

	_, err = nodesapp.RegisterService(ctx, a, logger.WithComponent("nodes"), manifests, authService, storesService, aliasService)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func getAPIKeys(ctx context.Context, cfg *csv.Config, logger log.Logger) (map[string]*authtypes.UserClaims, error) {
	if cfg != nil {
		apiKeyReader, err := csv.New(cfg)
		if err != nil {
			return nil, err
		}

		apikeyClaims, err := apiKeyReader.Load(ctx)
		if err != nil {
			return nil, err
		}

		logger.Info("API key authentication enabled")

		return apikeyClaims, nil
	}

	return nil, nil
}

func getJWTValidator(cfg *jose.Config, logger log.Logger) (jwt.Validator, error) {
	if cfg != nil {
		jwtValidator, err := jose.New(cfg)
		if err != nil {
			return nil, err
		}

		logger.Info("JWT authentication enabled")

		return jwtValidator, nil
	}

	return nil, nil
}

func getRootCAs(ctx context.Context, cfg *tls.Config, logger log.Logger) (*x509.CertPool, error) {
	if cfg != nil {
		tlsReader, err := tls.New(cfg)
		if err != nil {
			return nil, err
		}

		rootCAs, err := tlsReader.Load(ctx)
		if err != nil {
			return nil, err
		}

		logger.Info("TLS authentication enabled")

		return rootCAs, nil
	}

	return nil, nil
}

func getManifests(ctx context.Context, cfg *manifestreader.Config) (map[string][]entities2.Manifest, error) {
	manifestReader, err := manifestreader.New(cfg)
	if err != nil {
		return nil, err
	}

	return manifestReader.Load(ctx)
}
