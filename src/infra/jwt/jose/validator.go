package jose

import (
	"context"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/validate/josev2"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	jwtinfra "github.com/consensys/quorum-key-manager/src/infra/jwt"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Validator struct {
	validator *josev2.Validator
}

var _ jwtinfra.Validator = &Validator{}

func New(cfg *Config) (*Validator, error) {
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

	return &Validator{validator: validator}, nil
}

func (v *Validator) ValidateToken(ctx context.Context, token string) (*entities.UserClaims, error) {
	userCtx, err := v.validator.ValidateToken(ctx, token)
	if err != nil {
		// There is no fine-grained handling of the error provided from the package
		return nil, err
	}

	return &entities.UserClaims{
		Subject: userCtx.(*josev2.UserContext).Claims.Subject,
		Scope:   userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).Scope,
		Roles:   userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).Roles,
	}, nil
}
