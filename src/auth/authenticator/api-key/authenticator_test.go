package apikey

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	"golang.org/x/crypto/sha3"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	BobAPIKey   = "bobAPIKey"
	AliceAPIKey = "aliceAPIKey"
)

var userAliceClaims = UserClaims{
	UserName: "TenantOne|Alice",
	Claims:   []string{"guest", "admin", "read:key", "write:key"},
}

var userBobClaims = UserClaims{
	UserName: "Bob",
	Claims:   []string{"signer", "reader", "read:secret", "write:secret"},
}

func TestAuthenticatorApiKey_sh256Hasher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha256.Sum256([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha256.Sum256([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrAliceAPIKeyHash: userAliceClaims,
		hexStrBobAPIKeyHash:   userBobClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept api-key and extract ID successfully", func(t *testing.T) {
		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, "TenantOne", userInfo.Tenant)
		assert.Equal(t, []string{"guest", "admin"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:key", "write:key"}, userInfo.Permissions)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)

		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"signer", "reader"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:secret", "write:secret"}, userInfo.Permissions)
	})
}

func TestAuthenticatorApiKey_ChangingHashers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha512.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha512.Sum512([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha512.Sum512([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrAliceAPIKeyHash: userAliceClaims,
		hexStrBobAPIKeyHash:   userBobClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept api key and extract ID successfully", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"guest", "admin"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:key", "write:key"}, userInfo.Permissions)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)
		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"signer", "reader"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:secret", "write:secret"}, userInfo.Permissions)

	})
}

func TestAuthenticatorApiKey_base64encoder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha3.New224()
	b64Encoder := base64.URLEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha3.Sum224([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha3.Sum224([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrAliceAPIKeyHash: userAliceClaims,
		hexStrBobAPIKeyHash:   userBobClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept api key and extract ID successfully with Url encoding", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"guest", "admin"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:key", "write:key"}, userInfo.Permissions)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)

		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"signer", "reader"}, userInfo.Roles)
		assert.Equal(t, []types.Permission{"read:secret", "write:secret"}, userInfo.Permissions)
	})
}

func TestAuthenticatorApiKey_InvalidEncoder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha3.New224()
	b64Encoder := base64.URLEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := hasher.Sum([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash)

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrAliceAPIKeyHash: userAliceClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should reject api key with error", func(t *testing.T) {
		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, aliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.Error(t, err)
		assert.Nil(t, userInfo)
	})
}

func TestAuthenticatorApiKey_InvalidApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	anyWrongKey := "anyWrongKey"

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding
	b64AnyWrongAPIKey := b64Encoder.EncodeToString([]byte(anyWrongKey))

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := hasher.Sum([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash)

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrBobAPIKeyHash: userBobClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should reject api key and return error", func(t *testing.T) {
		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AnyWrongAPIKey))

		userInfo, err := auth.Authenticate(reqBob)

		require.Error(t, err)
		assert.Nil(t, userInfo)

	})
}

func TestAuthenticatorApiKey_EmptyApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := "aliceAPIKey"
	aliceAPIKeyHash := hasher.Sum([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash)

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserClaims{
		hexStrAliceAPIKeyHash: userAliceClaims,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should not reject api key and return error", func(t *testing.T) {
		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		userInfo, err := auth.Authenticate(reqAlice)
		require.NoError(t, err)
		assert.Nil(t, userInfo)
	})
}

func TestAuthenticatorApiKey_NilAuth(t *testing.T) {
	auth, _ := NewAuthenticator(&Config{})
	t.Run("should not instantiate when no config provided", func(t *testing.T) {
		assert.Nil(t, auth)
	})

}
