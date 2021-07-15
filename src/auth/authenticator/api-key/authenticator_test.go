package apikey

import (
	"crypto/sha256"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestAuthenticatorSameApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKey := "my-api-key"
	hasher := sha256.New()

	auth, _ := NewAuthenticator(&Config{ApiKeyFile: map[string]*UserNameAndGroups{},
		Hasher: hasher,
	})

	t.Run("should accept apikey and extract ID successfully", func(t *testing.T) {

		reqAlice := httptest.NewRequest("GET", "https://test.url", nil)
		reqAlice.Header.Add(ApiKeyHeader, string(apiKey))

		userInfo, err := auth.Authenticate(reqAlice)

		require.NoError(t, err)
		assert.Equal(t, "Alice", userInfo.Username)
		assert.Equal(t, []string{"Consensys"}, userInfo.Groups)

	})

}
