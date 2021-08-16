package entities

import "time"

type Metadata struct {
	Version     string
	Disabled    bool
	ExpireAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `pg:",soft_delete"`
	DestroyedAt time.Time
}
