package oidc

import (
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.MapClaims

	Username string   `json:"username"`
	Claims   []string `json:"claims"`

	cfg *ClaimsConfig
}

func (c *Claims) UnmarshalJSON(b []byte) error {
	// First Unmarshal JWT entries
	err := json.Unmarshal(b, &c.MapClaims)
	if err != nil {
		return err
	}

	// Second Unmarshal Orchestrate entries
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(b, &objmap)
	if err != nil {
		return err
	}

	if raw, ok := objmap[c.cfg.Subject]; ok {
		err = json.Unmarshal(*raw, &c.Username)
		if err != nil {
			return err
		}
	}

	if raw, ok := objmap[c.cfg.Scope]; ok {
		var claims string
		err = json.Unmarshal(*raw, &claims)
		if err != nil {
			return err
		}

		c.Claims = strings.Split(claims, " ")
	}

	return nil
}
