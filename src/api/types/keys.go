package types

import "time"

type CreateKeyRequest struct {
	ID               string            `json:"id" validate:"required" example:"my-key"`
	Curve            string            `json:"curve" validate:"required,isCurve" example:"secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type ImportKeyRequest struct {
	ID               string            `json:"id" validate:"required" example:"my-key"`
	Curve            string            `json:"curve" validate:"required,isCurve" example:"secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa"`
	PrivateKey       string            `json:"privateKey" validate:"required,isBase64" example:"0xfeee"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type SignPayloadRequest struct {
	Data string `json:"data" validate:"required,isBase64" example:"0xfeee"`
}

type KeyResponse struct {
	ID               string            `json:"id" example:"my-key"`
	PublicKey        string            `json:"publicKey" example:"0xfeee"`
	Curve            string            `json:"curve" example:"secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" example:"ecdsa"`
	Tags             map[string]string `json:"tags,omitempty"`
	Version          string            `json:"version" example:"1"`
	Disabled         bool              `json:"disabled" example:"false"`
	CreatedAt        time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt        time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	ExpireAt         time.Time         `json:"expireAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt        time.Time         `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt      time.Time         `json:"destroyedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}
