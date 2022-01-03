package authenticator

import (
	"context"
	"crypto/sha256"
	tls2 "crypto/tls"
	"crypto/x509"
	"fmt"
	"hash"
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
	hasher       hash.Hash
	rootCAs      *x509.CertPool
}

var _ auth.Authenticator = &Authenticator{}

func New(jwtValidator jwt.Validator, apiKeyClaims map[string]*entities.UserClaims, rootCAs *x509.CertPool, logger log.Logger) *Authenticator {
	return &Authenticator{
		jwtValidator: jwtValidator,
		apiKeyClaims: apiKeyClaims,
		rootCAs:      rootCAs,
		hasher:       sha256.New(),
		logger:       logger,
	}
}

func (authen *Authenticator) AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error) {
	if authen.jwtValidator == nil {
		errMessage := "jwt authentication method is not enabled"
		authen.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	authen.logger.Debug("extracting user info from jwt token")

	claims, err := authen.jwtValidator.ValidateToken(ctx, token)
	if err != nil {
		errMessage := "failed to validate jwt token"
		authen.logger.WithError(err).Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return authen.userInfoFromClaims(JWTAuthMode, claims), nil
}

func (authen *Authenticator) AuthenticateAPIKey(_ context.Context, apiKey []byte) (*entities.UserInfo, error) {
	if authen.apiKeyClaims == nil {
		errMessage := "api key authentication method is not enabled"
		authen.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	authen.logger.Debug("extracting user info from api key")

	authen.hasher.Reset()
	_, err := authen.hasher.Write(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to hash api key")
	}

	apiKeySha256 := fmt.Sprintf("%x", authen.hasher.Sum(nil))
	claims, ok := authen.apiKeyClaims[apiKeySha256]
	if !ok {
		errMessage := "invalid api key"
		authen.logger.Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return authen.userInfoFromClaims(APIKeyAuthMode, claims), nil
}

// AuthenticateTLS checks rootCAs and retrieve user info
func (authen Authenticator) AuthenticateTLS(_ context.Context, connState *tls2.ConnectionState) (*entities.UserInfo, error) {
	if authen.rootCAs == nil {
		errMessage := "tls authentication method is not enabled"
		authen.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	if !connState.HandshakeComplete {
		errMessage := "request must complete valid handshake"
		authen.logger.Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	err := tls.VerifyCertificateAuthority(connState.PeerCertificates, connState.ServerName, authen.rootCAs, true)
	if err != nil {
		errMessage := "invalid tls certificate"
		authen.logger.WithError(err).Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	// first array element is the leaf
	clientCert := connState.PeerCertificates[0]
	claims := &entities.UserClaims{
		Subject: clientCert.Subject.CommonName,
		Scope:   strings.Join(clientCert.Subject.OrganizationalUnit, " "),
		Roles:   strings.Join(clientCert.Subject.Organization, " "),
	}
	return authen.userInfoFromClaims(TLSAuthMode, claims), nil
}

func (authen *Authenticator) userInfoFromClaims(authMode string, claims *entities.UserClaims) *entities.UserInfo {
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

	authen.logger.Debug(
		"user info extracted successfully",
		"username", userInfo.Username,
		"tenant", userInfo.Tenant,
		"permissions", userInfo.Permissions,
		"roles", userInfo.Roles,
	)

	return userInfo
}
