package jose

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/common"
)

type Claims struct {
	Permissions     []string      `json:"permissions"`
	CustomClaims    *CustomClaims `json:"-"`
	customClaimPath string
	permissionsPath string
}

type CustomClaims struct {
	TenantID string `json:"tenant_id"`
}

func NewClaims(customClaimPath, permissionsPath string) *Claims {
	return &Claims{
		customClaimPath: customClaimPath,
		permissionsPath: permissionsPath,
	}
}

func (c *Claims) UnmarshalJSON(data []byte) error {
	var res map[string]interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	if c.customClaimPath != "" {
		c.CustomClaims = &CustomClaims{}
		if _, ok := res[c.customClaimPath]; ok {
			bClaims, _ := json.Marshal(res[c.customClaimPath])
			if err := json.Unmarshal(bClaims, &c.CustomClaims); err != nil || c.CustomClaims.TenantID == "" {
				return errors.New("invalid custom claims format")
			}
		} else {
			return errors.New("missing custom claims data")
		}
	}

	if c.permissionsPath != "" {
		if err := common.InterfaceToObject(res[c.permissionsPath], &c.Permissions); err != nil {
			return errors.New("invalid permission data type")
		}
	} else if res["scope"] != nil {
		c.Permissions = strings.Split(res["scope"].(string), " ")
	}

	return nil
}

func (c *Claims) Validate(_ context.Context) error {
	// TODO: Apply validation on custom claims if needed, currently no validation is needed
	return nil
}
