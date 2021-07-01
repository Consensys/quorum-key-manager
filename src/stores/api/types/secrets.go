package types

import "time"

type SetSecretRequest struct {
	Value string            `json:"value" validate:"required" example:"my-value"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type SecretResponse struct {
	ID          string            `json:"id" example:"my-secret"`
	Value       string            `json:"value" example:"my-value"`
	Tags        map[string]string `json:"tags,omitempty"`
	Version     string            `json:"version" example:"1"`
	Disabled    bool              `json:"disabled" example:"false"`
	CreatedAt   time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt   time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	ExpireAt    time.Time         `json:"expireAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt   time.Time         `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt time.Time         `json:"destroyedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}
