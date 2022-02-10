package jose

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

type Claims struct {
	CustomClaims    *CustomClaims `json:"-"`
	Scope           []string      `json:"scope"`
	customClaimPath string
}

type CustomClaims struct {
	TenantID    string   `json:"tenant_id"`
	Permissions []string `json:"permissions"`
}

func NewClaims(customClaimPath string) *Claims {
	return &Claims{
		customClaimPath: customClaimPath,
	}
}

func (c *Claims) UnmarshalJSON(data []byte) error {
	c.Scope = nil
	c.CustomClaims = nil

	var res map[string]interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	if c.customClaimPath != "" {
		c.CustomClaims = &CustomClaims{}
		if _, ok := res[c.customClaimPath]; ok {
			bClaims, _ := json.Marshal(res[c.customClaimPath])
			if err := json.Unmarshal(bClaims, &c.CustomClaims); err != nil {
				return errors.New("invalid custom claims format")
			}
			if c.CustomClaims.TenantID == "" {
				return errors.New("custom claims must include tenant_id")
			}
		} else {
			return errors.New("missing custom claims data")
		}
	}

	if res["scope"] != nil {
		if scopes, ok := res["scope"].(string); ok {
			c.Scope = strings.Split(scopes, " ")
		}
	}

	return nil
}

func (c *Claims) Validate(_ context.Context) error {
	// TODO: Apply validation on custom claims if needed, currently no validation is needed
	return nil
}
