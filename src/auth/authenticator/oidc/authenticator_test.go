package oidc

import (
	"crypto/x509"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/jwt"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/pkg/tls/testutils"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticator_RSAToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	claimsCfg := &ClaimsConfig{
		Subject: "sub",
		Scope:   "scope",
		Roles:   "qkm-user-roles",
	}

	cert, _ := certificate.X509KeyPair([]byte(testutils.RSACertPEM), []byte(testutils.RSAKeyPEM))
	generator, _ := jwt.NewTokenGenerator(cert.PrivateKey)
	auth, _ := NewAuthenticator(&Config{
		Certificates: []*x509.Certificate{cert.Leaf},
		Claims:       claimsCfg,
	})

	t.Run("should accept token and extract claims successfully", func(t *testing.T) {
		claims := []string{"read:key", "write:key"}
		roles := []string{"operator", "signer"}
		token, _ := generator.GenerateAccessToken(map[string]interface{}{
			claimsCfg.Subject: "tenant|username",
			claimsCfg.Scope:   strings.Join(claims, " "),
			claimsCfg.Roles:   strings.Join(roles, ","),
		}, time.Second)

		req := httptest.NewRequest("GET", "http://test.url", nil)
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", BearerSchema, token))
		userInfo, err := auth.Authenticate(req)
		require.NoError(t, err)
		assert.Equal(t, "username", userInfo.Username)
		assert.Equal(t, "tenant", userInfo.Tenant)
		assert.Equal(t, []string{"operator", "signer"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:key", "write:key"}, userInfo.Permissions)
	})

	t.Run("should reject request for invalid token", func(t *testing.T) {
		token := "invalid-auth-token"
		req := httptest.NewRequest("GET", "http://test.url", nil)
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", BearerSchema, token))
		userInfo, err := auth.Authenticate(req)
		assert.Error(t, err)
		assert.Empty(t, userInfo)
	})

	t.Run("should ignore authenticator for missing token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://test.url", nil)
		userInfo, err := auth.Authenticate(req)
		assert.Nil(t, err)
		assert.Empty(t, userInfo)
	})
}
