package jose

import (
	"testing"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Parser(t *testing.T) {
	v := Validator{}

	t.Run("should parse standard token successfully", func(t *testing.T) {
		tokenClaims := &validator.ValidatedClaims{
			CustomClaims: nil,
			RegisteredClaims: validator.RegisteredClaims{
				Subject: "tenant_id",
			},
		}
		c, err := v.ParseClaims(tokenClaims)
		assert.NoError(t, err)
		assert.Equal(t, "tenant_id", c.Tenant)
	})

	t.Run("should parse token with permissions successfully", func(t *testing.T) {
		tokenClaims := &validator.ValidatedClaims{
			CustomClaims: &Claims{
				Permissions: []string{"read:*", "*:keys"},
			},
			RegisteredClaims: validator.RegisteredClaims{
				Subject: "tenant_id",
			},
		}
		c, err := v.ParseClaims(tokenClaims)
		assert.NoError(t, err)
		assert.Equal(t, "tenant_id", c.Tenant)
		assert.Equal(t, []string{"read:*", "*:keys"}, c.Permissions)
	})

	t.Run("should parse token with permissions and custom claims successfully", func(t *testing.T) {
		tokenClaims := &validator.ValidatedClaims{
			CustomClaims: &Claims{
				CustomClaims: &CustomClaims{
					TenantID: "tenant_id_2",
				},
				Permissions: []string{"read:*", "*:keys"},
			},
			RegisteredClaims: validator.RegisteredClaims{
				Subject: "tenant_id",
			},
		}
		c, err := v.ParseClaims(tokenClaims)
		assert.NoError(t, err)
		assert.Equal(t, "tenant_id_2", c.Tenant)
		assert.Equal(t, []string{"read:*", "*:keys"}, c.Permissions)
	})

	t.Run("should fail if invalid token is passed", func(t *testing.T) {
		tokenClaims := validator.ValidatedClaims{}
		_, err := v.ParseClaims(tokenClaims)
		assert.Error(t, err)
	})
}
