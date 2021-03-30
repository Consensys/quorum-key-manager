package entities

import "time"

type Metadata struct {
	Version     int
	Disabled    bool
	ExpireAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	DestroyedAt time.Time
}
