package authenticator

import (
	"context"
	tls2 "crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/tls"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/jwt"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

const (
	APIKeyAuthMode = "apikey"
	JWTAuthMode    = "jwt"
	TLSAuthMode    = "tls"
)

type Authenticator struct {
	logger       log.Logger
	jwtValidator jwt.Validator
	apiKeyClaims map[string]*entities.UserClaims
	rootCAs      *x509.CertPool
}

var _ auth.Authenticator = &Authenticator{}

func New(jwtValidator jwt.Validator, apiKeyClaims map[string]*entities.UserClaims, rootCAs *x509.CertPool, logger log.Logger) *Authenticator {
	return &Authenticator{
		jwtValidator: jwtValidator,
		apiKeyClaims: apiKeyClaims,
		rootCAs:      rootCAs,
		logger:       logger,
	}
}

func (auth *Authenticator) AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from jwt token")

	claims, err := auth.jwtValidator.ValidateToken(ctx, token)
	if err != nil {
		errMessage := "failed to validate jwt token"
		auth.logger.WithError(err).Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return auth.userInfoFromClaims(JWTAuthMode, claims), nil
}

func (auth *Authenticator) AuthenticateAPIKey(_ context.Context, apiKey string) (*entities.UserInfo, error) {
	auth.logger.Debug("extracting user info from api key")

	claims, ok := auth.apiKeyClaims[apiKey]
	if !ok {
		errMessage := "api key not found"
		auth.logger.Warn(errMessage, "api_key_hash", apiKey)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return auth.userInfoFromClaims(APIKeyAuthMode, claims), nil
}

// AuthenticateTLS checks rootCAs and retrieve user info
func (auth Authenticator) AuthenticateTLS(_ context.Context, connState *tls2.ConnectionState) (*entities.UserInfo, error) {
	if !connState.HandshakeComplete {
		errMessage := "request must complete valid handshake"
		auth.logger.Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	err := tls.VerifyCertificateAuthority(connState.PeerCertificates, connState.ServerName, auth.rootCAs, true)
	if err != nil {
		errMessage := "invalid tls certificate"
		auth.logger.WithError(err).Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	// first array element is the leaf
	clientCert := connState.PeerCertificates[0]
	claims := &entities.UserClaims{
		Subject: clientCert.Subject.CommonName,
		Scope:   strings.Join(clientCert.Subject.OrganizationalUnit, " "),
		Roles:   strings.Join(clientCert.Subject.Organization, " "),
	}
	return auth.userInfoFromClaims(TLSAuthMode, claims), nil
}

func (auth *Authenticator) userInfoFromClaims(authMode string, claims *entities.UserClaims) *entities.UserInfo {
	userInfo := &entities.UserInfo{AuthMode: authMode}

	// If more than one element in subject, then the username has been specified
	subject := strings.Split(claims.Subject, "|")
	if len(subject) > 1 {
		userInfo.Username = subject[1]
	}
	userInfo.Tenant = subject[0]

	for _, permission := range strings.Fields(claims.Scope) {
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
		userInfo.Roles = strings.Fields(claims.Roles)
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
