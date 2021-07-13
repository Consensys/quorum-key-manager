package authenticator

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

var authenticatedGroup = "system:authenticated"

// Middleware synchronize authentication
type Middleware struct {
	authenticator Authenticator
	logger        log.Logger
}

func NewMiddleware(logger log.Logger, authenticators ...Authenticator) *Middleware {
	return &Middleware{
		authenticator: First(authenticators...),
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

	mid.logger.With("groups", info.Groups).Debug("request successfully authenticated")

	ctx = WithUserContext(ctx, NewUserContext(info))

	// Serve next
	next.ServeHTTP(rw, req.WithContext(ctx))
}

func OnError(w http.ResponseWriter, _ *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}
