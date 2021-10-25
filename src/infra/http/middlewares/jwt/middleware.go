package jwt

import (
	"context"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/auth0/go-jwt-middleware/validate/josev2"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http/middlewares"
	"github.com/consensys/quorum-key-manager/src/infra/http/utils"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"net/url"
	"time"
)

const authMode = "JWT"

type Middleware struct {
	validator *josev2.Validator
}

var _ httpinfra.Middleware = &Middleware{}

func New(cfg *Config) (*Middleware, error) {
	issuerURL, err := url.Parse(cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	validator, err := josev2.New(
		josev2.NewCachingJWKSProvider(*issuerURL, cfg.CacheTTL).KeyFunc,
		jose.RS256,
		josev2.WithCustomClaims(func() josev2.CustomClaims { return &CustomClaims{} }),
		josev2.WithExpectedClaims(func() jwt.Expected {
			return jwt.Expected{
				Audience: cfg.Audience,
				Time:     time.Now(),
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Middleware{validator: validator}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return jwtmiddleware.New(
		m.validateToken,
		jwtmiddleware.WithCredentialsOptional(true),
		jwtmiddleware.WithErrorHandler(parseErrorResponse),
	).CheckJWT(next)
}

func (m *Middleware) validateToken(ctx context.Context, token string) (interface{}, error) {
	userCtx, err := m.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	subject := userCtx.(*josev2.UserContext).Claims.Subject
	scope := userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).Scope
	roles := userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).Roles

	userInfo := &entities.UserInfo{AuthMode: authMode}
	userInfo.Username, userInfo.Tenant = utils.ExtractUsernameAndTenant(subject)
	userInfo.Permissions = utils.ExtractPermissions(scope)
	if roles != "" {
		userInfo.Roles = utils.ExtractRoles(roles)
	}

	return userInfo, nil
}
