package entities

import "time"

// Secret
type Secret struct {
	Value       string
	Disabled    bool
	Recovery    *Recovery
	Tags        map[string]string
	Version     int
	ExpireAt    time.Time
	CreatedAt   time.Time
	DeletedAt   time.Time
	DestroyedAt time.Time
}
