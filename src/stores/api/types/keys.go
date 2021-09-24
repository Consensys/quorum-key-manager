package types

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type CreateKeyRequest struct {
	Curve            string            `json:"curve" validate:"required,isCurve" example:"secp256k1" enums:"babyjubjub,secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa" enums:"ecdsa,eddsa"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type ImportKeyRequest struct {
	Curve            string            `json:"curve" validate:"required,isCurve" example:"secp256k1" enums:"babyjubjub,secp256k1"`
	SigningAlgorithm string            `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa" enums:"ecdsa,eddsa"`
	PrivateKey       []byte            `json:"privateKey" validate:"required" example:"bXkgc2lnbmVkIG1lc3NhZ2U=" swaggertype:"string"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type UpdateKeyRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type SignBase64PayloadRequest struct {
	Data []byte `json:"data" validate:"required" example:"bXkgc2lnbmVkIG1lc3NhZ2U=" swaggertype:"string"`
}

type VerifyKeySignatureRequest struct {
	Data             []byte `json:"data" validate:"required" example:"bXkgc2lnbmVkIG1lc3NhZ2U=" swaggertype:"string"`
	Signature        []byte `json:"signature" validate:"required" example:"tjThYhKSFSKKvsR8Pji6EJ+FYAcf8TNUdAQnM7MSwZEEaPvFhpr1SuGpX5uOcYUrb3pBA8cLk8xcbKtvZ56qWA==" swaggertype:"string"`
	Curve            string `json:"curve" validate:"required,isCurve" example:"secp256k1" enums:"babyjubjub,secp256k1" swaggertype:"string"`
	SigningAlgorithm string `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa" enums:"ecdsa,eddsa"`
	PublicKey        []byte `json:"publicKey" validate:"required" example:"Cjix/fS3WdqKGKabagBNYwcClan5aImoFpnjSF0cqJs=" swaggertype:"string"`
}

type KeyResponse struct {
	ID               string               `json:"id" example:"my-key"`
	PublicKey        string               `json:"publicKey" example:"Cjix/fS3WdqKGKabagBNYwcClan5aImoFpnjSF0cqJs=" swaggertype:"string"`
	Curve            string               `json:"curve" example:"secp256k1"`
	SigningAlgorithm string               `json:"signingAlgorithm" example:"ecdsa"`
	Tags             map[string]string    `json:"tags,omitempty"`
	Annotations      *entities.Annotation `json:"annotations,omitempty"`
	Disabled         bool                 `json:"disabled" example:"false"`
	CreatedAt        time.Time            `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt        time.Time            `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt        *time.Time           `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}
