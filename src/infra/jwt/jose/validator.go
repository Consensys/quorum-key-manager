package jose

import (
	"errors"
	"net/url"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
)

type Validator struct {
	*validator.Validator
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
			return NewClaims(cfg.CustomClaimPath)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Validator{v}, nil
}

func (v *Validator) ParseClaims(tokenClaims interface{}) (*entities.UserClaims, error) {
	claims, ok := tokenClaims.(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	userClaims := &entities.UserClaims{}
	if qkmUserClaims, ok := v.qkmCustomClaimsExist(claims); ok {
		userClaims.Tenant = qkmUserClaims.TenantID
		userClaims.Permissions = qkmUserClaims.Permissions
	} else {
		userClaims.Tenant = claims.RegisteredClaims.Subject
		if scopeClaims, ok := v.scopeClaimsExist(claims); ok {
			userClaims.Permissions = scopeClaims
		}
	}

	return userClaims, nil
}

func (v *Validator) qkmCustomClaimsExist(claims *validator.ValidatedClaims) (*CustomClaims, bool) {
	if claims.CustomClaims == nil {
		return nil, false
	}

	if qkmUserClaims := claims.CustomClaims.(*Claims).CustomClaims; qkmUserClaims != nil {
		return qkmUserClaims, true
	}

	return nil, false
}

func (v *Validator) scopeClaimsExist(claims *validator.ValidatedClaims) ([]string, bool) {
	if claims.CustomClaims == nil {
		return nil, false
	}

	if scopeUserClaims := claims.CustomClaims.(*Claims).Scope; scopeUserClaims != nil {
		return scopeUserClaims, true
	}

	return nil, false
}
