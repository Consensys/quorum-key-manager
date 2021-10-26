package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http"
	"net/http"
	"strings"
)

const BasicSchema = "Basic"

type AuthHandler struct {
	authenticator auth.Authenticator
}

func New(authenticator auth.Authenticator) *AuthHandler {
	return &AuthHandler{
		authenticator: authenticator,
	}
}

func (h *AuthHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// JWT Token
		jwtToken, err := jwtmiddleware.AuthHeaderTokenExtractor(r)
		if err != nil {
			httpinfra.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		}

		if jwtToken != "" {
			userInfo, err := h.authenticator.AuthenticateJWT(r.Context(), jwtToken)
			if err != nil {
				httpinfra.WriteHTTPErrorResponse(rw, err)
			}

			next.ServeHTTP(rw, r.Clone(WithUserInfo(ctx, userInfo)))
		}

		// API key
		apiKey, err := extractApiKey(r)
		if err != nil {
			httpinfra.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
			return
		}

		if apiKey != nil {
			userInfo, err := h.authenticator.AuthenticateAPIKey(r.Context(), apiKey)
			if err != nil {
				httpinfra.WriteHTTPErrorResponse(rw, err)
			}

			next.ServeHTTP(rw, r.Clone(WithUserInfo(ctx, userInfo)))
		}

		// TLS
		// TODO: Implement TLS authenticator

		// Anonymous user if no authentication method has succeeded
		next.ServeHTTP(rw, r.Clone(WithUserInfo(ctx, entities.NewAnonymousUser())))
	})
}

func extractApiKey(r *http.Request) ([]byte, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return nil, nil
	}

	if len(authHeader) <= len(BasicSchema) || !strings.EqualFold(authHeader[:len(BasicSchema)], BasicSchema) {
		return nil, fmt.Errorf("api key was not provided in Authorization header")
	}

	b64EncodedAPIKey := authHeader[len(BasicSchema)+1:]
	decodedAPIKey, err := base64.StdEncoding.DecodeString(b64EncodedAPIKey)
	if err != nil {
		return nil, err
	}

	return decodedAPIKey, nil
}
