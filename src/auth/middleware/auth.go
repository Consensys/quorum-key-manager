package authmiddleware

import (
	"context"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/log"
	"github.com/consensys/quorum-key-manager/src/auth/authorization"
	"github.com/consensys/quorum-key-manager/src/auth/manager"
	"github.com/consensys/quorum-key-manager/src/auth/middleware/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

var authenticatedGroup = "system:authenticated"

type Middleware struct {
	authenticator authenticator.Authenticator
	policyMngr    manager.Manager
}

func New(auth authenticator.Authenticator, policyMngr manager.Manager) *Middleware {
	return &Middleware{
		authenticator: auth,
		policyMngr:    policyMngr,
	}
}

func (mid *Middleware) Then(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		mid.ServeHTTP(rw, req, h)
	})
}

func (mid *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	ctx := req.Context()
	logger := log.FromContext(ctx)

	// Authenticate request
	info, err := mid.authenticator.Authenticate(req)
	if err != nil {
		OnError(rw, req, err)
		return
	}

	if info != nil {
		// If authentication succeeded then sets the system:authenticated group
		info.Groups = append(info.Groups, authenticatedGroup)
	} else {
		// If no authentication then sets info to anonymous user
		info = types.AnonymousUser
	}

	// Create policy resolver for user info
	resolver, err := mid.authorizationResolver(ctx, info)
	if err != nil {
		logger.WithError(err).Errorf("could not create policy resolver")
		OnError(rw, req, err)
		return
	}
	ctx = authorization.WithResolver(ctx, resolver)

	// Create request context and sets user info
	reqCtx := types.NewRequestContext(req)
	reqCtx.UserInfo = info

	ctx = types.WithRequestContext(ctx, reqCtx)

	// Serve next
	next.ServeHTTP(rw, req.WithContext(ctx))
}

func (mid *Middleware) authorizationResolver(ctx context.Context, info *types.UserInfo) (*authorization.Resolver, error) {
	logger := log.FromContext(ctx)

	// Retrieve policies associated to user info
	var policies []*types.Policy
	for _, groupName := range info.Groups {
		group, err := mid.policyMngr.Group(ctx, groupName)
		if err != nil {
			logger.WithError(err).WithField("group", groupName).Debugf("could not load group")
			continue
		}

		for _, policyName := range group.Policies {
			policy, err := mid.policyMngr.Policy(ctx, policyName)
			if err != nil {
				logger.WithError(err).WithField("policy", groupName).Debugf("could not load policy")
				continue
			}
			policies = append(policies, policy)
		}
	}

	// Create resolver
	return authorization.NewResolver(policies)
}

func OnError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}
