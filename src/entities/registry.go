package entities

import "time"

type AliasRegistry struct {
	Name      string
	Aliases   []Alias
	CreatedAt time.Time
	UpdatedAt time.Time
}
