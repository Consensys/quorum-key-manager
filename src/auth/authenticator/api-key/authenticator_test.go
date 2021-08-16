package apikey

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/sha3"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	BobAPIKey   = "bobAPIKey"
	AliceAPIKey = "aliceAPIKey"
)

func TestAuthenticatorCorrectAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha256.Sum256([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])
	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha256.Sum256([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))

	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	userBobAndGroups := UserNameAndGroups{
		UserName: "Bob",
		Groups:   []string{"g3", "g1"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
		hexStrBobAPIKeyHash:   userBobAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept apikey and extract ID successfully", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"g1", "g2"}, userInfo.Groups)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)

		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"g3", "g1"}, userInfo.Groups)

	})
}

func TestAuthenticatorCorrectAPIKeyWithChangingHashers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha512.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha512.Sum512([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])
	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha512.Sum512([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))

	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	userBobAndGroups := UserNameAndGroups{
		UserName: "Bob",
		Groups:   []string{"g3", "g1"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
		hexStrBobAPIKeyHash:   userBobAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept apikey and extract ID successfully", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"g1", "g2"}, userInfo.Groups)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)

		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"g3", "g1"}, userInfo.Groups)

	})
}

func TestAuthenticatorCorrectAPIKeyWithChangingEncoder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha3.New224()
	b64Encoder := base64.URLEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := sha3.Sum224([]byte(aliceAPIKey))
	b64AliceAPIKey := b64Encoder.EncodeToString([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash[:])

	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := sha3.Sum224([]byte(bobAPIKey))
	b64BobAPIKey := b64Encoder.EncodeToString([]byte(bobAPIKey))

	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash[:])

	userBobAndGroups := UserNameAndGroups{
		UserName: "Bob",
		Groups:   []string{"g3", "g1"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
		hexStrBobAPIKeyHash:   userBobAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should accept apikey and extract ID successfully with Url encoding", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"g1", "g2"}, userInfo.Groups)

		reqBob := httptest.NewRequest("GET", "https://test.url", nil)
		reqBob.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64BobAPIKey))

		userInfo, err = auth.Authenticate(reqBob)

		require.NoError(t, err)
		assert.Equal(t, "Bob", userInfo.Username)
		assert.Equal(t, []string{"g3", "g1"}, userInfo.Groups)

	})
}

func TestAuthenticatorWrongEncodingAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha3.New224()
	b64Encoder := base64.URLEncoding

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := hasher.Sum([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash)
	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should reject apikey with error", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, aliceAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.Error(t, err)
		assert.Nil(t, userInfo)
	})
}

func TestAuthenticatorWrongAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	anyWrongKey := "anyWrongKey"

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding

	b64AnyWrongAPIKey := b64Encoder.EncodeToString([]byte(anyWrongKey))

	aliceAPIKey := AliceAPIKey
	aliceAPIKeyHash := hasher.Sum([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash)
	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	bobAPIKey := BobAPIKey
	bobAPIKeyHash := hasher.Sum([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash)

	userBobAndGroups := UserNameAndGroups{
		UserName: "Bob",
		Groups:   []string{"g3", "g1"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
		hexStrBobAPIKeyHash:   userBobAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should reject apikey and return error", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add("Authorization", fmt.Sprintf("%s %s", BasicSchema, b64AnyWrongAPIKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.Error(t, err)
		assert.Nil(t, userInfo)

	})
}

func TestAuthenticatorNoAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hasher := sha256.New()
	b64Encoder := base64.StdEncoding

	aliceAPIKey := "aliceAPIKey"
	aliceAPIKeyHash := hasher.Sum([]byte(aliceAPIKey))
	hexStrAliceAPIKeyHash := hex.EncodeToString(aliceAPIKeyHash)
	userAliceAndGroups := UserNameAndGroups{
		UserName: "Alice",
		Groups:   []string{"g1", "g2"},
	}

	bobAPIKey := "bobAPIKey"
	bobAPIKeyHash := hasher.Sum([]byte(bobAPIKey))
	hexStrBobAPIKeyHash := hex.EncodeToString(bobAPIKeyHash)

	userBobAndGroups := UserNameAndGroups{
		UserName: "Bob",
		Groups:   []string{"g3", "g1"},
	}

	auth, _ := NewAuthenticator(&Config{APIKeyFile: map[string]UserNameAndGroups{
		hexStrAliceAPIKeyHash: userAliceAndGroups,
		hexStrBobAPIKeyHash:   userBobAndGroups,
	},
		Hasher:     &hasher,
		B64Encoder: b64Encoder,
	})

	t.Run("should not reject missing apikey", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Nil(t, userInfo)

	})
}

func TestNilAuthenticator(t *testing.T) {

	auth, _ := NewAuthenticator(&Config{})

	t.Run("should not instantiate when no config provided", func(t *testing.T) {
		assert.Nil(t, auth)
	})

}
