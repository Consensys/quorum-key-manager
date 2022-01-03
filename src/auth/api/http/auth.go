package http

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http"
)

const BasicSchema = "basic"
const BearerSchema = "bearer"

type Auth struct {
	authenticator auth.Authenticator
}

func NewAuth(authenticator auth.Authenticator) *Auth {
	return &Auth{
		authenticator: authenticator,
	}
}

func (m *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")

		// If Auth header is provided, try JWT or API key
		if authHeader != "" {
			authHeaderParts := strings.Fields(authHeader)
			if len(authHeaderParts) != 2 {
				httpinfra.WriteHTTPErrorResponse(rw, errors.InvalidFormatError("malformed authorization header"))
				return
			}
			authSchema := authHeaderParts[0]
			authValue := authHeaderParts[1]

			switch strings.ToLower(authSchema) {
			case BearerSchema:
				userInfo, err := m.authenticator.AuthenticateJWT(r.Context(), authValue)
				if err != nil {
					httpinfra.WriteHTTPErrorResponse(rw, err)
					return
				}

				next.ServeHTTP(rw, r.WithContext(WithUserInfo(ctx, userInfo)))
				return
			case BasicSchema:
				apiKey, err := base64.StdEncoding.DecodeString(authValue)
				if err != nil {
					httpinfra.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
					return
				}

				userInfo, err := m.authenticator.AuthenticateAPIKey(r.Context(), apiKey)
				if err != nil {
					httpinfra.WriteHTTPErrorResponse(rw, err)
					return
				}

				next.ServeHTTP(rw, r.WithContext(WithUserInfo(ctx, userInfo)))
				return
			default:
				httpinfra.WriteHTTPErrorResponse(rw, errors.InvalidFormatError("unsupported authorization schema %s", authSchema))
				return
			}
		}

		// if no Authorization header, try TLS, otherwise anonymous
		if r.TLS != nil && r.TLS.PeerCertificates != nil && len(r.TLS.PeerCertificates) > 0 {
			userInfo, err := m.authenticator.AuthenticateTLS(r.Context(), r.TLS)
			if err != nil {
				httpinfra.WriteHTTPErrorResponse(rw, err)
				return
			}

			next.ServeHTTP(rw, r.WithContext(WithUserInfo(ctx, userInfo)))
			return
		}

		// Anonymous user if no authentication method has succeeded
		next.ServeHTTP(rw, r.WithContext(WithUserInfo(ctx, entities.NewAnonymousUser())))
	})
}
