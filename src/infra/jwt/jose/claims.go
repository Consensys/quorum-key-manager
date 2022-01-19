package jose

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

type Claims struct {
	Roles           []string      `json:"roles"`
	Permissions     []string      `json:"permissions"`
	CustomClaims    *CustomClaims `json:"-"`
	customClaimPath string
	permissionsPath string
	rolesPath       string
}

type CustomClaims struct {
	TenantID string `json:"tenant_id"`
}

func NewClaims(customClaimPath, permissionsPath, rolesPath string) *Claims {
	return &Claims{
		customClaimPath: customClaimPath,
		permissionsPath: permissionsPath,
		rolesPath:       rolesPath,
	}
}

func (c *Claims) UnmarshalJSON(data []byte) error {
	var res map[string]interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	if c.customClaimPath != "" {
		if _, ok := res[c.customClaimPath]; ok {
			c.CustomClaims = &CustomClaims{}
			bClaims, _ := json.Marshal(res[c.customClaimPath])
			if err := json.Unmarshal(bClaims, &c.CustomClaims); err != nil {
				return errors.New("invalid custom claims format")
			}
		} else {
			return errors.New("missing custom claims data")
		}
	}

	if c.permissionsPath != "" {
		c.Permissions = res[c.permissionsPath].([]string)
	} else {
		c.Permissions = strings.Split(res["scope"].(string), "")
	}

	if c.rolesPath != "" {
		c.Roles = res[c.rolesPath].([]string)
	}

	return nil
}

func (claims *Claims) Validate(_ context.Context) error {
	// TODO: Apply validation on custom claims if needed, currently no validation is needed
	return nil
}
