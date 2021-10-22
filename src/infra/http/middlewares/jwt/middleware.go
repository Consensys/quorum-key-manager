package jwt

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/auth0/go-jwt-middleware/validate/josev2"
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http/middlewares"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"net/url"
)

type Middleware struct {
	middleware *jwtmiddleware.JWTMiddleware
}

var _ httpinfra.Middleware = &Middleware{}

func NewMiddleware(cfg *Config) (*Middleware, error) {
	issuerURL, err := url.Parse(cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	validator, err := josev2.New(
		josev2.NewCachingJWKSProvider(*issuerURL, cfg.CacheTTL).KeyFunc,
		jose.RS256,
	)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		middleware: jwtmiddleware.New(validator.ValidateToken, jwtmiddleware.WithCredentialsOptional(true)),
	}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return m.middleware.CheckJWT(next)
}
