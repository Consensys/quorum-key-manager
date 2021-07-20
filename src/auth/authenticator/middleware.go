package authenticator

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

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
		mid.onError(rw, req, err)
		return
	}

	if info != nil {
		// If authentication succeeded then sets the system:authenticated group
		info.Groups = append(types.AuthenticatedUser.Groups, info.Groups...)
	} else {
		// If no authentication then sets info to anonymous user
		info = types.AnonymousUser
	}

	mid.logger.With("groups", info.Groups).Debug("request successfully authenticated")

	ctx = WithUserContext(ctx, NewUserContext(info))

	// Serve next
	next.ServeHTTP(rw, req.WithContext(ctx))
}

func (mid *Middleware) onError(w http.ResponseWriter, _ *http.Request, err error) {
	errMsg := "unauthorized request"
	mid.logger.Error(errMsg, "err", err.Error())
	http.Error(w, errMsg, http.StatusUnauthorized)
}
