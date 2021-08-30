package authenticator

import (
	"net/http"

	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

// Middleware synchronize authentication
type Middleware struct {
	authenticator Authenticator
	authEnabled   bool
	logger        log.Logger
}

func NewMiddleware(logger log.Logger, authenticators ...Authenticator) *Middleware {
	return &Middleware{
		authenticator: First(authenticators...),
		authEnabled:   len(authenticators) > 0,
		logger:        logger,
	}
}

func (mid *Middleware) Then(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		mid.ServeHTTP(rw, req, h)
	})
}

func (mid *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	ctx := req.Context()
	if !mid.authEnabled {
		ctx = WithUserContext(ctx, NewUserContext(types.WildcardUser))
		next.ServeHTTP(rw, req.WithContext(ctx))
		return
	}

	// Authenticate request
	info, err := mid.authenticator.Authenticate(req)
	if err != nil {
		mid.logger.WithError(err).Error("unauthorized request")
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	if info != nil {
		// If authentication succeeded then sets the system:authenticated group
		mid.logger.With("username", info.Username).
			With("tenant", info.Tenant).
			With("roles", info.Roles).
			With("permissions", info.Permissions).
			Debug("request successfully authenticated")
	} else {
		// If no authentication then sets info to anonymous user
		info = types.AnonymousUser
		mid.logger.With("username", info.Username).
			With("roles", info.Roles).
			With("permissions", info.Permissions).
			Debug("anonymous request received")
	}

	ctx = WithUserContext(ctx, NewUserContext(info))

	// Serve next
	next.ServeHTTP(rw, req.WithContext(ctx))
}
