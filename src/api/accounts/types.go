package accountsapi

import (
	"time"
)

type Metadata struct {
	Version int `json:"version"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	PurgeAt   time.Time `json:"purgeAt"`
}

type CreateAccountRequest struct {
	StoreName string            `json:"name"`
	Enabled   bool              `json:"enabled"`
	ExpireAt  time.Time         `json:"expireAt"`
	Tags      map[string]string `json:"tags"`
}

type CreateAccountResponse struct {
	Addr string `json:"address"`

	Metadata *Metadata
}
