package app

import (
	"context"
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/pkg/app"
	service "github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/api/http"
	"github.com/consensys/quorum-key-manager/src/auth/api/manifest"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/service/roles"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/justinas/alice"
)

func RegisterService(
	ctx context.Context,
	a *app.App,
	logger log.Logger,
	manifests map[string][]entities2.Manifest,
	jwtValidator jwt.Validator,
	apikeyClaims map[string]*authtypes.UserClaims,
	rootCAs *x509.CertPool,
) (*roles.Roles, error) {
	// Business layer
	// TODO: Create authorizator service here

	var authenticatorService *authenticator.Authenticator
	if jwtValidator != nil || apikeyClaims != nil || rootCAs != nil {
		authenticatorService = authenticator.New(jwtValidator, apikeyClaims, rootCAs, logger)
	}

	rolesService := roles.New(logger)

	// Service layer
	manifestRolesHandler := manifest.NewRolesHandler(rolesService) // Manifest reading is synchronous, similar to a config file
	for _, mnf := range manifests[entities2.RoleKind] {
		err := manifestRolesHandler.Create(ctx, mnf.Specs)
		if err != nil {
			return nil, err
		}
	}

	compositionMiddleware, err := createMiddlewares(logger, authenticatorService)
	if err != nil {
		return nil, err
	}

	err = a.SetMiddleware(compositionMiddleware.Then)
	if err != nil {
		return nil, err
	}

	return rolesService, nil
}

func createMiddlewares(logger log.Logger, authenticator service.Authenticator) (alice.Chain, error) {
	var authMiddleware alice.Constructor

	if authenticator != nil {
		authMiddleware = http.NewAuth(authenticator).Middleware
	} else {
		logger.Warn("No authentication method enabled")
		authMiddleware = http.WildcardMiddleware
	}

	return alice.New(
		http.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		authMiddleware,
	), nil
}
