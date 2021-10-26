package authenticator

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	apikey "github.com/consensys/quorum-key-manager/src/infra/api-key"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"strings"
)

const (
	APIKeyAuthMode = "apikey"
	JWTAuthMode    = "jwt"
)

type Authenticator struct {
	logger       log.Logger
	jwtValidator jwt.Validator
	apikeyReader apikey.Reader
}

var _ auth.Authenticator = &Authenticator{}

func New(jwtValidator jwt.Validator, apikeyReader apikey.Reader, logger log.Logger) *Authenticator {
	return &Authenticator{
		jwtValidator: jwtValidator,
		apikeyReader: apikeyReader,
		logger:       logger,
	}
}

func (auth *Authenticator) AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from jwt token")

	claims, err := auth.jwtValidator.ValidateToken(ctx, token)
	if err != nil {
		auth.logger.WithError(err).Error("failed to validate jwt token")
		return nil, err
	}

	return auth.userInfoFromClaims(JWTAuthMode, claims), nil
}

func (auth *Authenticator) AuthenticateAPIKey(ctx context.Context, apiKey []byte) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from api key")

	claims, err := auth.apikeyReader.Get(ctx, apiKey)
	if err != nil {
		auth.logger.WithError(err).Error("failed to validate api key")
		return nil, err
	}

	return auth.userInfoFromClaims(APIKeyAuthMode, claims), nil
}

func (auth *Authenticator) userInfoFromClaims(authMode string, claims *entities.UserClaims) *entities.UserInfo {
	userInfo := &entities.UserInfo{AuthMode: authMode}

	// If more than one element in subject, then the username has been specified
	subject := strings.Split(claims.Subject, "|")
	if len(subject) > 1 {
		userInfo.Username = subject[1]
	}
	userInfo.Tenant = subject[0]

	for _, permission := range strings.Split(claims.Scope, " ") {
		if !strings.Contains(permission, ":") {
			// Ignore invalid permissions
			continue
		}

		if strings.Contains(permission, "*") {
			userInfo.Permissions = append(userInfo.Permissions, entities.ListWildcardPermission(permission)...)
		} else {
			userInfo.Permissions = append(userInfo.Permissions, entities.Permission(permission))
		}
	}

	if claims.Roles != "" {
		userInfo.Roles = strings.Split(claims.Roles, " ")
	}

	auth.logger.Debug(
		"user info extracted successfully",
		"username", userInfo.Username,
		"tenant", userInfo.Tenant,
		"permissions", userInfo.Permissions,
		"roles", userInfo.Roles,
	)

	return userInfo
}
