package jose

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaims_standard(t *testing.T) {
	c := NewClaims("", "")

	t.Run("should parse token successfully", func(t *testing.T) {
		token := map[string]interface{}{
			"scope": "read:* write:ethereum",
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, c.Permissions, strings.Split(token["scope"].(string), " "))
		assert.Nil(t, c.CustomClaims)
	})

	t.Run("should not fail to parse token with no scope", func(t *testing.T) {
		token := map[string]interface{}{}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.NoError(t, err)
	})

	t.Run("should fail if invalid token data", func(t *testing.T) {
		err := c.UnmarshalJSON([]byte("invalid data"))
		assert.Error(t, err)
	})
}

func TestClaims_customClaims(t *testing.T) {
	customClaimPath := "my.custom.claim"
	c := NewClaims(customClaimPath, "")

	t.Run("should parse token with custom claims successfully", func(t *testing.T) {
		token := map[string]interface{}{
			"scope":         "read:* write:ethereum",
			customClaimPath: map[string]string{"tenant_id": "tenantID"},
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, c.Permissions, strings.Split(token["scope"].(string), " "))
		assert.Equal(t, c.CustomClaims.TenantID, "tenantID")
	})

	t.Run("should fail to parse token with invalid custom claims type", func(t *testing.T) {
		token := map[string]interface{}{
			customClaimPath: map[string]int{"tenant_id": 12},
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.Error(t, err)
	})

	t.Run("should fail to parse token with invalid custom claims format ", func(t *testing.T) {
		token := map[string]interface{}{
			customClaimPath: map[string]string{"xx": "tenantID"},
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.Error(t, err)
	})

	t.Run("should fail if missing custom claims", func(t *testing.T) {
		token := map[string]interface{}{}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}

func TestClaims_customPermissions(t *testing.T) {
	customPermissionPath := "permissions"

	c := NewClaims("", customPermissionPath)

	t.Run("should parse token with custom permissions successfully", func(t *testing.T) {
		token := map[string]interface{}{
			"scope":              "read:* write:ethereum",
			customPermissionPath: []string{"read:*", "*:keys"},
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, c.Permissions, token[customPermissionPath])
	})

	t.Run("should fail to parse token with invalid permission format", func(t *testing.T) {
		token := map[string]interface{}{
			customPermissionPath: "read:* *:keys",
		}

		data, _ := json.Marshal(token)
		err := c.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}

func TestClaims_validate(t *testing.T) {
	c := NewClaims("", "")
	assert.Equal(t, nil, c.Validate(context.Background()))
}
