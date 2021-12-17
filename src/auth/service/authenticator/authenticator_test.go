package authenticator

import (
	"context"
	"crypto/sha256"
	tls2 "crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"

	mock2 "github.com/consensys/quorum-key-manager/src/infra/log/mock"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities/testdata"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/stretchr/testify/suite"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bobAPIKey   = "bobAPIKey"
	aliceAPIKey = "aliceAPIKey"
)

type authenticatorTestSuite struct {
	suite.Suite
	mockJWTValidator *mock.MockValidator
	userClaims       map[string]*entities.UserClaims
	aliceCert        *x509.Certificate
	eveCert          *x509.Certificate
	logger           *mock2.MockLogger
	auth             *Authenticator
}

func TestAuthenticator(t *testing.T) {
	s := new(authenticatorTestSuite)
	suite.Run(t, s)
}

func (s *authenticatorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	// User claims
	aliceClaims := testdata.FakeUserClaims()
	bobClaims := testdata.FakeUserClaims()
	bobClaims.Scope = "*:*"

	aliceSha256 := fmt.Sprintf("%x", sha256.Sum256([]byte(aliceAPIKey)))
	bobSha256 := fmt.Sprintf("%x", sha256.Sum256([]byte(bobAPIKey)))

	s.userClaims = map[string]*entities.UserClaims{
		aliceSha256: aliceClaims, // base64 of alice key
		bobSha256:   bobClaims,   // base64 of bob key
	}

	// TLS certs
	aliceCert, err := certificate.X509KeyPair([]byte(testdata.TLSClientAliceCert), []byte(testdata.TLSAuthKey))
	require.NoError(s.T(), err)
	eveCert, err := certificate.X509KeyPair([]byte(testdata.TLSClientEveCert), []byte(testdata.TLSAuthKeyEve))
	require.NoError(s.T(), err)

	s.aliceCert = aliceCert.Leaf
	s.eveCert = eveCert.Leaf

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(s.aliceCert)
	caCertPool.AddCert(s.eveCert)

	s.mockJWTValidator = mock.NewMockValidator(ctrl)
	s.logger = testutils2.NewMockLogger(ctrl)

	s.auth = New(s.mockJWTValidator, s.userClaims, caCertPool, s.logger)
}

func (s *authenticatorTestSuite) TestAuthenticateJWT() {
	ctx := context.Background()
	token := "myToken"

	s.Run("should authenticate a jwt token successfully", func() {
		s.mockJWTValidator.EXPECT().ValidateToken(ctx, token).Return(testdata.FakeUserClaims(), nil)

		userInfo, err := s.auth.AuthenticateJWT(ctx, token)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), "Alice", userInfo.Username)
		assert.Equal(s.T(), "TenantOne", userInfo.Tenant)
		assert.Equal(s.T(), []string{"guest", "admin"}, userInfo.Roles)
		assert.Equal(s.T(), []entities.Permission{"read:key", "write:key"}, userInfo.Permissions)
		assert.Equal(s.T(), JWTAuthMode, userInfo.AuthMode)
	})

	s.Run("should authenticate a jwt token successfully with wildcard permissions", func() {
		userClaims := testdata.FakeUserClaims()
		userClaims.Scope = "*:*"
		s.mockJWTValidator.EXPECT().ValidateToken(ctx, token).Return(userClaims, nil)

		userInfo, err := s.auth.AuthenticateJWT(ctx, token)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), entities.NewWildcardUser().Permissions, userInfo.Permissions)
	})

	s.Run("should return UnauthorizedError if the token fails validation", func() {
		s.mockJWTValidator.EXPECT().ValidateToken(ctx, token).Return(testdata.FakeUserClaims(), fmt.Errorf("error"))

		userInfo, err := s.auth.AuthenticateJWT(ctx, token)

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})

	s.Run("should return UnauthorizedError if the authentication method is not enabled", func() {
		auth := New(nil, nil, nil, s.logger)

		userInfo, err := auth.AuthenticateJWT(ctx, token)

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})
}

func (s *authenticatorTestSuite) TestAuthenticateAPIKey() {
	ctx := context.Background()

	s.Run("should authenticate with api key successfully", func() {
		userInfo, err := s.auth.AuthenticateAPIKey(ctx, []byte(aliceAPIKey))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), "Alice", userInfo.Username)
		assert.Equal(s.T(), "TenantOne", userInfo.Tenant)
		assert.Equal(s.T(), []string{"guest", "admin"}, userInfo.Roles)
		assert.Equal(s.T(), []entities.Permission{"read:key", "write:key"}, userInfo.Permissions)
		assert.Equal(s.T(), APIKeyAuthMode, userInfo.AuthMode)
	})

	s.Run("should authenticate an api key successfully with wildcard permissions", func() {
		userInfo, err := s.auth.AuthenticateAPIKey(ctx, []byte(bobAPIKey))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), entities.NewWildcardUser().Permissions, userInfo.Permissions)
	})

	s.Run("should return UnauthorizedError if api key is not found", func() {
		userInfo, err := s.auth.AuthenticateAPIKey(ctx, []byte("invalid-key"))

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})

	s.Run("should return UnauthorizedError if the authentication method is not enabled", func() {
		auth := New(nil, nil, nil, s.logger)

		userInfo, err := auth.AuthenticateAPIKey(ctx, []byte(aliceAPIKey))

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})
}

func (s *authenticatorTestSuite) TestAuthenticateTLS() {
	ctx := context.Background()

	s.Run("should authenticate with TLS successfully", func() {
		connState := &tls2.ConnectionState{
			PeerCertificates: []*x509.Certificate{s.aliceCert},
		}
		connState.HandshakeComplete = true

		userInfo, err := s.auth.AuthenticateTLS(ctx, connState)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), "", userInfo.Username)
		assert.Equal(s.T(), "alice", userInfo.Tenant)
		assert.Equal(s.T(), []string{"admin", "signer"}, userInfo.Roles)
		assert.Equal(s.T(), []entities.Permission{"read:accounts", "delete:secrets"}, userInfo.Permissions)
		assert.Equal(s.T(), TLSAuthMode, userInfo.AuthMode)
	})

	s.Run("should authenticate an api key successfully with wildcard permissions", func() {
		connState := &tls2.ConnectionState{
			PeerCertificates: []*x509.Certificate{s.eveCert},
		}
		connState.HandshakeComplete = true

		userInfo, err := s.auth.AuthenticateTLS(ctx, connState)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), []entities.Permission{"read:ethereum", "write:ethereum", "delete:ethereum", "destroy:ethereum", "sign:ethereum", "encrypt:ethereum"}, userInfo.Permissions)
	})

	s.Run("should return UnauthorizedError if tls has not handshaked", func() {
		connState := &tls2.ConnectionState{
			PeerCertificates: []*x509.Certificate{s.eveCert},
		}
		connState.HandshakeComplete = false

		userInfo, err := s.auth.AuthenticateTLS(ctx, connState)

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})

	s.Run("should return UnauthorizedError if the authentication method is not enabled", func() {
		connState := &tls2.ConnectionState{
			PeerCertificates: []*x509.Certificate{s.eveCert},
		}
		connState.HandshakeComplete = false

		auth := New(nil, nil, nil, s.logger)

		userInfo, err := auth.AuthenticateTLS(ctx, connState)

		require.Nil(s.T(), userInfo)
		assert.True(s.T(), errors.IsUnauthorizedError(err))
	})
}
