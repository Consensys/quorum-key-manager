package types

import (
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ECRecoverRequest struct {
	Data      hexutil.Bytes `json:"data" validate:"required" example:"0xfeaeee..." swaggertype:"string"`
	Signature hexutil.Bytes `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
}

type VerifyRequest struct {
	Data      hexutil.Bytes  `json:"data" validate:"required" example:"0xfeaeee..." swaggertype:"string"`
	Signature hexutil.Bytes  `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
	Address   common.Address `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
}

type VerifyTypedDataRequest struct {
	TypedData types.SignTypedDataRequest `json:"data" validate:"required"`
	Signature hexutil.Bytes              `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
	Address   common.Address             `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
}

type VerifyKeySignatureRequest struct {
	Data             []byte `json:"data" validate:"required" example:"bXkgc2lnbmVkIG1lc3NhZ2U=" swaggertype:"string"`
	Signature        []byte `json:"signature" validate:"required" example:"tjThYhKSFSKKvsR8Pji6EJ+FYAcf8TNUdAQnM7MSwZEEaPvFhpr1SuGpX5uOcYUrb3pBA8cLk8xcbKtvZ56qWA==" swaggertype:"string"`
	Curve            string `json:"curve" validate:"required,isCurve" example:"secp256k1" enums:"babyjubjub,secp256k1" swaggertype:"string"`
	SigningAlgorithm string `json:"signingAlgorithm" validate:"required,isSigningAlgorithm" example:"ecdsa" enums:"ecdsa,eddsa"`
	PublicKey        []byte `json:"publicKey" validate:"required" example:"Cjix/fS3WdqKGKabagBNYwcClan5aImoFpnjSF0cqJs=" swaggertype:"string"`
}
