package types

import "time"

type CreateSecretRequest struct {
	ID    string            `json:"id" example:"my-privateKey" validate:"required"`
	Value string            `json:"value" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249" validate:"required"`
	Tags  map[string]string `json:"tags"`
}

type Secret struct {
	ID       string            `json:"id" example:"my-privateKey"`
	Value    string            `json:"value" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"`
	Metadata *Metadata         `json:"metadata"`
	Tags     map[string]string `json:"tags"`
}

type Metadata struct {
	Version     int       `json:"version" example:"my-privateKey"`
	Disabled    bool      `json:"disabled" example:"false"`
	ExpireAt    time.Time `json:"expireAt" example:"2020-07-09T12:35:42.115395Z"`
	CreatedAt   time.Time `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt   time.Time `json:"deletedAt" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt time.Time `json:"destroyedAt" example:"2020-07-09T12:35:42.115395Z"`
}
