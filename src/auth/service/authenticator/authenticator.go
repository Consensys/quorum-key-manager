package authenticator

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"hash"
	"strings"
)

const (
	APIKeyAuthMode = "apikey"
	JWTAuthMode    = "jwt"
)

type Authenticator struct {
	logger       log.Logger
	jwtvalidator jwt.Validator
	APIKeyFile   map[string]UserClaims
	Hasher       *hash.Hash
	B64Encoder   *base64.Encoding
}

var _ auth.Authenticator = &Authenticator{}

func New(jwtvalidator jwt.Validator, logger log.Logger) *Authenticator {
	return &Authenticator{
		jwtvalidator: jwtvalidator,
		logger:       logger,
	}
}

func (auth *Authenticator) AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from jwt token")

	claims, err := auth.jwtvalidator.ValidateToken(ctx, token)
	if err != nil {
		errMessage := "failed to validate jwt token"
		auth.logger.WithError(err).Error(errMessage)
		return nil, err
	}

	return auth.userInfoFromClaims(JWTAuthMode, claims), nil
}

func (auth *Authenticator) AuthenticateAPIKey(ctx context.Context, apiKey []byte) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from api key file")

	h := *auth.Hasher
	h.Reset()
	_, err := h.Write(apiKey)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}
	clientAPIKeyHash := h.Sum(nil)

	strClientHash := hex.EncodeToString(clientAPIKeyHash)
	claims, ok := auth.APIKeyFile[strClientHash]
	if !ok {
		return nil, errors.UnauthorizedError("invalid api-key")
	}

	userInfo := &entities.UserInfo{
		AuthMode:    AuthMode,
		Roles:       []string{},
		Permissions: []entities.Permission{},
	}

	userInfo.Username, userInfo.Tenant = ExtractUsernameAndTenant(claims.UserName)
	userInfo.Permissions = ExtractPermissionsArr(claims.Permissions)
	userInfo.Roles = claims.Roles

	auth.logger.Info(
		"user info extracted from api key successfully",
		"username", userInfo.Username,
		"tenant", userInfo.Tenant,
		"permissions", userInfo.Permissions,
		"roles", userInfo.Roles,
	)

	return userInfo, nil
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

	auth.logger.Info(
		"user info extracted successfully",
		"username", userInfo.Username,
		"tenant", userInfo.Tenant,
		"permissions", userInfo.Permissions,
		"roles", userInfo.Roles,
	)

	return userInfo
}
