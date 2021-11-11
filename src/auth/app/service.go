package app

import (
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/pkg/app"
	service "github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/api/middlewares"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/manager"
	"github.com/consensys/quorum-key-manager/src/auth/service/authenticator"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/justinas/alice"
)

func RegisterService(
	a *app.App,
	logger log.Logger,
	manifests []entities2.Manifest,
	jwtValidator jwt.Validator,
	apikeyClaims map[string]*entities.UserClaims,
	rootCAs *x509.CertPool,
) error {
	// Business layer
	// TODO: Create authorizator service here

	var authenticatorService service.Authenticator
	if jwtValidator != nil || apikeyClaims != nil || rootCAs != nil {
		authenticatorService = authenticator.New(jwtValidator, apikeyClaims, rootCAs, logger)
	}

	// Service layer
	compositionMiddleware, err := createMiddlewares(logger, authenticatorService)
	if err != nil {
		return err
	}

	err = a.SetMiddleware(compositionMiddleware.Then)
	if err != nil {
		return err
	}

	// TODO: Remove manager
	policyMngr := manager.New(manifests, logger)
	return a.RegisterService(policyMngr)
}

func createMiddlewares(logger log.Logger, authenticator service.Authenticator) (alice.Chain, error) {
	var authMiddleware alice.Constructor

	if authenticator != nil {
		authMiddleware = middlewares.NewAuth(authenticator).Middleware
	} else {
		logger.Warn("No authentication method enabled")
		authMiddleware = middlewares.WildcardMiddleware
	}

	return alice.New(
		middlewares.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		authMiddleware,
	), nil
}
