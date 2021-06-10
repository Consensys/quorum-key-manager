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
	PrivateKey       []byte            `json:"privateKey" validate:"required" example:"BFVSFJhqUh9DQJwcayNtsWdD2..." swaggertype:"string"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type UpdateKeyRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type SignBase64PayloadRequest struct {
	Data []byte `json:"data" validate:"required" example:"BFVSFJhqUh9DQJwcayNtsWdD2..." swaggertype:"string"`
}

type VerifyKeySignatureRequest struct {
	Data             []byte `json:"data" validate:"required" example:"BFVSFJhqUh9DQJwcayNtsWdD2..."`
	Signature        []byte `json:"signature" validate:"required" example:"BFVSFJhqUh9DQJwcayNtsWdD2..."`
	Curve            string `json:"curve" validate:"required,isCurve" example:"secp256k1"`
	SigningAlgorithm string `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa"`
	PublicKey        []byte `json:"publicKey" validate:"required" example:"BFVSFJhqUh9DQJwcayNtsWdD2..."`
}

type KeyResponse struct {
	ID               string            `json:"id" example:"my-key"`
	PublicKey        string            `json:"publicKey" example:"BFVSFJhqUh9DQJwcayNtsWdD2..."`
	Curve            string            `json:"curve" example:"secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" example:"ecdsa"`
	Tags             map[string]string `json:"tags,omitempty"`
	Disabled         bool              `json:"disabled" example:"false"`
	CreatedAt        time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt        time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	ExpireAt         time.Time         `json:"expireAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt        time.Time         `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt      time.Time         `json:"destroyedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}
