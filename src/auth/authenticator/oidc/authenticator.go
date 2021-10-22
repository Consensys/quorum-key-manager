package oidc

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/utils"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"net/http"
)

const AuthMode = "JWT"
const BearerSchema = "Bearer"

type Authenticator struct {
	jwtChecker *JWTChecker
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) *Authenticator {
	return &Authenticator{jwtChecker: NewJWTChecker(cfg.IssuerURL, cfg.Claims, false)}
}

func (a Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// Extract Access Token from context
	token, ok := extractToken(BearerSchema, req.Header.Get("Authorization"))
	// In case of no credentials are sent we authenticate with Anonymous user
	if !ok {
		return nil, nil
	}

	jwtData, err := a.jwtChecker.Check(req.Context(), token)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}

	userInfo := &types.UserInfo{
		AuthMode: AuthMode,
	}

	userInfo.Username, userInfo.Tenant = utils.ExtractUsernameAndTenant(jwtData.Subject)
	userInfo.Permissions = utils.ExtractPermissions(jwtData.Scope)
	rolesClaim := a.jwtChecker.claimsCfg.Roles
	if jwtData.MapClaims[rolesClaim] != nil {
		userInfo.Roles = utils.ExtractRoles(jwtData.MapClaims[rolesClaim].(string))
	}
	return userInfo, nil
}
