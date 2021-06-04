package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CreateEth1AccountRequest struct {
	ID   string            `json:"id" validate:"required" example:"my-account"`
	Tags map[string]string `json:"tags,omitempty"`
}

type ImportEth1AccountRequest struct {
	ID         string            `json:"id" validate:"required" example:"my-account"`
	PrivateKey hexutil.Bytes     `json:"privateKey" validate:"required" example:"0xfeee"`
	Tags       map[string]string `json:"tags,omitempty"`
}

type UpdateEth1AccountRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type SignHexPayloadRequest struct {
	Data hexutil.Bytes `json:"data" validate:"required" example:"0xfeee"`
}

type SignTypedDataRequest struct {
	DomainSeparator DomainSeparator        `json:"domainSeparator" validate:"required"`
	Types           map[string][]Type      `json:"types" validate:"required"`
	Message         map[string]interface{} `json:"message" validate:"required"`
	MessageType     string                 `json:"messageType" validate:"required" example:"Mail"`
}

type DomainSeparator struct {
	Name              string `json:"name" validate:"required" example:"MyDApp"`
	Version           string `json:"version" validate:"required" example:"v1.0.0"`
	ChainID           int64  `json:"chainID" validate:"required" example:"1"`
	VerifyingContract string `json:"verifyingContract,omitempty" validate:"omitempty,isHexAddress" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Salt              string `json:"salt,omitempty" validate:"omitempty" example:"some-random-string"`
}

type Type struct {
	Name string `json:"name" validate:"required" example:"fieldName"`
	Type string `json:"type" validate:"required" example:"string"`
}

type SignETHTransactionRequest struct {
	Nonce    hexutil.Uint64  `json:"nonce" example:"0x1"`
	To       *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Value    hexutil.Big     `json:"value,omitempty" example:"0xfeaeae"`
	GasPrice hexutil.Big     `json:"gasPrice" validate:"required" example:"0x0"`
	GasLimit hexutil.Uint64  `json:"gasLimit" validate:"required" example:"0x5208"`
	Data     hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..."`
	ChainID  hexutil.Big     `json:"chainID" validate:"required" example:"0x1 (mainnet)"`
}

type SignQuorumPrivateTransactionRequest struct {
	Nonce    hexutil.Uint64  `json:"nonce" example:"0x1"`
	To       *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Value    hexutil.Big     `json:"value,omitempty" example:"0x1"`
	GasPrice hexutil.Big     `json:"gasPrice" validate:"required" example:"0x0"`
	GasLimit hexutil.Uint64  `json:"gasLimit" validate:"required" example:"0x5208"`
	Data     hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..."`
}

type SignEEATransactionRequest struct {
	Nonce          hexutil.Uint64  `json:"nonce" example:"0x1"`
	To             *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Data           hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..."`
	ChainID        hexutil.Big     `json:"chainID" validate:"required" example:"0x1 (mainnet)"`
	PrivateFrom    string          `json:"privateFrom" validate:"required,base64,required_with=PrivateFor PrivacyGroupID" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string        `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string          `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

type ECRecoverRequest struct {
	Data      hexutil.Bytes `json:"data" validate:"required" example:"0xfeaeee..."`
	Signature hexutil.Bytes `json:"signature" validate:"required" example:"0x6019a3c8..."`
}

type VerifyEth1SignatureRequest struct {
	Data      hexutil.Bytes  `json:"data" validate:"required" example:"0xfeaeee..."`
	Signature hexutil.Bytes  `json:"signature" validate:"required" example:"0x6019a3c8..."`
	Address   common.Address `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
}

type VerifyTypedDataRequest struct {
	TypedData SignTypedDataRequest `json:"data" validate:"required"`
	Signature hexutil.Bytes        `json:"signature" validate:"required" example:"0x6019a3c8..."`
	Address   common.Address       `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
}

// Represents the ethereum accounts 
// swagger:model Eth1AccountResponse
type Eth1AccountResponse struct {
	PublicKey           hexutil.Bytes     `json:"publicKey" example:"0xfeee"`
	CompressedPublicKey hexutil.Bytes     `json:"compressedPublicKey" example:"0xfeee"`
	CreatedAt           time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt           time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	ExpireAt            time.Time         `json:"expireAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt           time.Time         `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	DestroyedAt         time.Time         `json:"destroyedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	ID                  string            `json:"id" example:"my-key"`
	Tags                map[string]string `json:"tags,omitempty"`
	Address             common.Address    `json:"address" example:"0xfeee"`
	Disabled            bool              `json:"disabled" example:"false"`
}
