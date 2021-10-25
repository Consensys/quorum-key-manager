package jwt

import (
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	errorspkg "github.com/consensys/quorum-key-manager/pkg/errors"
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http"
	"net/http"
)

func parseErrorResponse(rw http.ResponseWriter, _ *http.Request, err error) {
	switch {
	case errors.Is(err, jwtmiddleware.ErrJWTMissing):
		httpinfra.WriteHTTPErrorResponse(rw, errorspkg.InvalidFormatError(err.Error()))
	case errors.Is(err, jwtmiddleware.ErrJWTInvalid):
		httpinfra.WriteHTTPErrorResponse(rw, errorspkg.UnauthorizedError(err.Error()))
	default:
		httpinfra.WriteHTTPErrorResponse(rw, err)
	}
}
