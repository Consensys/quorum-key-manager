package entities

import "time"

type AliasRegistry struct {
	Name           string
	Aliases        []Alias
	AllowedTenants []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
