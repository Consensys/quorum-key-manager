package types

import (
	"time"
)

type CreateEth1AccountRequest struct {
	ID string `json:"id" validate:"required" example:"my-id"`
}

type Eth1Response struct {
	ID                  string    `json:"id" example:"my-key"`
	Address             string    `json:"address" example:"0xfeee"`
	PublicKey           string    `json:"publicKey" example:"0xfeee"`
	CompressedPublicKey string    `json:"compressedPublicKey" example:"0xfeee"`
	Disabled            bool      `json:"disabled" example:"false"`
	CreatedAt           time.Time `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt           time.Time `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	ExpireAt            time.Time `json:"expireAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt           time.Time `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt         time.Time `json:"destroyedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}
