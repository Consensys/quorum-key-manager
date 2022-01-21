package jose

import (
	"context"
	"net/url"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
)

type Validator struct {
	validator *validator.Validator
}

var _ jwt.Validator = &Validator{}

func New(cfg *Config) (*Validator, error) {
	issuerURL, err := url.Parse(cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	v, err := validator.New(
		jwks.NewCachingProvider(issuerURL, cfg.CacheTTL).KeyFunc,
		validator.RS256,
		issuerURL.String(),
		cfg.Audience,
		validator.WithCustomClaims(func() validator.CustomClaims {
			return NewClaims(cfg.CustomClaimPath, cfg.PermissionClaimPath)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Validator{validator: v}, nil
}

func (v *Validator) ValidateToken(ctx context.Context, token string) (*entities.UserClaims, error) {
	userCtx, err := v.validator.ValidateToken(ctx, token)
	if err != nil {
		// There is no fine-grained handling of the error provided from the package
		return nil, err
	}

	claims := userCtx.(*validator.ValidatedClaims)
	userClaims := &entities.UserClaims{}

	if claims.CustomClaims != nil {
		userClaims.Permissions = claims.CustomClaims.(*Claims).Permissions
		if qkmUserClaims := claims.CustomClaims.(*Claims).CustomClaims; qkmUserClaims != nil {
			userClaims.Tenant = qkmUserClaims.TenantID
		} else {
			userClaims.Tenant = claims.RegisteredClaims.Subject
		}
	} else {
		userClaims.Tenant = claims.RegisteredClaims.Subject
	}

	return userClaims, nil
}
