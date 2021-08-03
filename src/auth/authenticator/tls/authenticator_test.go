package tls

import (
	"crypto/tls"
	"crypto/x509"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/pkg/tls/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatorSameCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aliceCert, _ := certificate.X509KeyPair([]byte(testutils.TLSClientAliceCert), []byte(testutils.TLSAuthKey))

	auth, _ := NewAuthenticator(&Config{})

	t.Run("should accept cert and extract ID successfully", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.TLS = &tls.ConnectionState{}
		reqAlice.TLS.PeerCertificates = make([]*x509.Certificate, 1)
		reqAlice.TLS.PeerCertificates[0] = aliceCert.Leaf

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"Consensys"}, userInfo.Groups)

	})

}

func TestAuthenticatorDifferentCert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aliceCert, _ := certificate.X509KeyPair([]byte(testutils.TLSClientAliceCert), []byte(testutils.TLSAuthKey))
	eveCert, _ := certificate.X509KeyPair([]byte(testutils.TLSClientEveCert), []byte(testutils.TLSAuthKey))

	assert.NotEqualValues(t, aliceCert, eveCert)

	auth, _ := NewAuthenticator(&Config{})

	t.Run("should NOT reject cert and leave ID empty", func(t *testing.T) {

		reqEve := httptest.NewRequest("GET", "https://test.url", nil)
		reqEve.TLS = &tls.ConnectionState{}
		reqEve.TLS.PeerCertificates = nil

		userInfo, err := auth.Authenticate(reqEve)

		require.NoError(t, err)
		assert.Nil(t, userInfo)
	})

}

func TestAuthenticator(t *testing.T) {

	auth, _ := NewAuthenticator(&Config{})

	t.Run("should instantiate new authenticator", func(t *testing.T) {
		assert.NotNil(t, auth)
	})

}
