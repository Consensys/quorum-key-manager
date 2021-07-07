package auth

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/log"
	authmanager "github.com/consensys/quorum-key-manager/src/auth/manager"
	authmiddleware "github.com/consensys/quorum-key-manager/src/auth/middleware"
	authenticator "github.com/consensys/quorum-key-manager/src/auth/middleware/authenticator"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load manifests service
	m := new(manifestsmanager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Create and register the stores service
	policyMngr := authmanager.New(*m, logger)
	err = a.RegisterService(policyMngr)
	if err != nil {
		return err
	}

	return nil
}

func Middleware(a *app.App, logger log.Logger) (func(http.Handler) http.Handler, error) {
	// Load configuration
	cfg := new(Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Load policy manager service
	policyMngr := new(authmanager.Manager)
	err = a.Service(policyMngr)
	if err != nil {
		return nil, err
	}

	// Create middleware
	mid := authmiddleware.New(
		authenticator.First(
		// TODO: pass each authenticator implementation based on config
		),
		*policyMngr,
		logger,
	)

	return mid.Then, nil
}
