package app

import (
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/pkg/app"
	service "github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/api/http"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/service/roles"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/justinas/alice"
)

func RegisterService(
	a *app.App,
	logger log.Logger,
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
	err := a.SetMiddleware(createMiddlewares(logger, authenticatorService).Then)
	if err != nil {
		return nil, err
	}

	return rolesService, nil
}

func createMiddlewares(logger log.Logger, authenticator service.Authenticator) alice.Chain {
	return alice.New(
		http.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		http.NewAuth(authenticator).Middleware,
	)
}
