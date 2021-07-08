package oicd

import (
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.MapClaims

	Username string   `json:"username"`
	Groups   []string `json:"groups"`

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

	if raw, ok := objmap[c.cfg.Username]; ok {
		err = json.Unmarshal(*raw, &c.Username)
		if err != nil {
			return err
		}
	}

	if raw, ok := objmap[c.cfg.Group]; ok {
		var groups string
		err = json.Unmarshal(*raw, &groups)
		if err != nil {
			return err
		}

		c.Groups = strings.Split(groups, ",")
	}

	return nil
}
