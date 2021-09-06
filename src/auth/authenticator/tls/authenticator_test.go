package tls

import (
	"crypto/tls"
	"crypto/x509"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/pkg/tls/testutils"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatorSameCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aliceCert, err := certificate.X509KeyPair([]byte(testutils.TLSClientAliceCert), []byte(testutils.TLSAuthKey))
	require.NoError(t, err)
	eveCert, err := certificate.X509KeyPair([]byte(testutils.TLSClientEveCert), []byte(testutils.TLSAuthKey))
	require.NoError(t, err)

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(aliceCert.Leaf)
	caCertPool.AddCert(eveCert.Leaf)
	auth, _ := NewAuthenticator(&Config{CAs: caCertPool})

	t.Run("should accept cert and extract username and roles successfully", func(t *testing.T) {
		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.TLS = &tls.ConnectionState{}
		reqAlice.TLS.PeerCertificates = make([]*x509.Certificate, 1)
		reqAlice.TLS.PeerCertificates[0] = aliceCert.Leaf
		reqAlice.TLS.HandshakeComplete = true

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "alice", userInfo.Username)
		assert.Equal(t, []string{"admin", "signer"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:accounts", "delete:secrets"}, userInfo.Permissions)

	})

	t.Run("should accept cert and extract username|tenant and roles|permissions successfully", func(t *testing.T) {
		reqEve := httptest.NewRequest("GET", "https://test.url", nil)
		reqEve.TLS = &tls.ConnectionState{}
		reqEve.TLS.PeerCertificates = make([]*x509.Certificate, 1)
		reqEve.TLS.PeerCertificates[0] = eveCert.Leaf
		reqEve.TLS.HandshakeComplete = true

		userInfo, err := auth.Authenticate(reqEve)

		require.NoError(t, err)
		assert.Equal(t, "eve", userInfo.Username)
		assert.Equal(t, "auth0", userInfo.Tenant)
		assert.Equal(t, []string{"signer"}, userInfo.Roles)
		assert.Equal(t, types.ListWildcardPermission("*:ethereum"), userInfo.Permissions)
	})
}

func TestAuthenticatorDifferentCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aliceCert, _ := certificate.X509KeyPair([]byte(testutils.TLSClientAliceCert), []byte(testutils.TLSAuthKey))
	eveCert, _ := certificate.X509KeyPair([]byte(testutils.TLSClientEveCert), []byte(testutils.TLSAuthKey))
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(aliceCert.Leaf)
	caCertPool.AddCert(eveCert.Leaf)
	auth, _ := NewAuthenticator(&Config{CAs: caCertPool})

	t.Run("should NOT reject cert and leave ID empty", func(t *testing.T) {

		reqEve := httptest.NewRequest("GET", "https://test.url", nil)
		reqEve.TLS = &tls.ConnectionState{}
		reqEve.TLS.PeerCertificates = nil

		userInfo, err := auth.Authenticate(reqEve)

		require.NoError(t, err)
		assert.Nil(t, userInfo)
	})

}

func TestEmptyAuthenticator(t *testing.T) {
	auth, _ := NewAuthenticator(&Config{})
	t.Run("should not instantiate new authenticator", func(t *testing.T) {
		assert.Nil(t, auth)
	})

}
