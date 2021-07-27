package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CreateEth1AccountRequest struct {
	KeyID string            `json:"keyId" example:"my-key-account"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type ImportEth1AccountRequest struct {
	KeyID      string            `json:"keyId" example:"my-imported-key-account"`
	PrivateKey hexutil.Bytes     `json:"privateKey" validate:"required" example:"0x56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E" swaggertype:"string"`
	Tags       map[string]string `json:"tags,omitempty"`
}

type UpdateEth1AccountRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type SignHexPayloadRequest struct {
	// required to be hex value
	Data hexutil.Bytes `json:"data" validate:"required" example:"0xfeee" swaggertype:"string"`
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
	Nonce    hexutil.Uint64  `json:"nonce" example:"0x1" swaggertype:"string"`
	To       *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
	Value    hexutil.Big     `json:"value,omitempty" example:"0xfeaeae" swaggertype:"string"`
	GasPrice hexutil.Big     `json:"gasPrice" validate:"required" example:"0x0" swaggertype:"string"`
	GasLimit hexutil.Uint64  `json:"gasLimit" validate:"required" example:"0x5208" swaggertype:"string"`
	Data     hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..." swaggertype:"string"`
	ChainID  hexutil.Big     `json:"chainID" validate:"required" example:"0x1 (mainnet)" swaggertype:"string"`
}

type SignQuorumPrivateTransactionRequest struct {
	Nonce    hexutil.Uint64  `json:"nonce" example:"0x1" swaggertype:"string"`
	To       *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
	Value    hexutil.Big     `json:"value,omitempty" example:"0x1" swaggertype:"string"`
	GasPrice hexutil.Big     `json:"gasPrice" validate:"required" example:"0x0" swaggertype:"string"`
	GasLimit hexutil.Uint64  `json:"gasLimit" validate:"required" example:"0x5208" swaggertype:"string"`
	Data     hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..." swaggertype:"string"`
}

type SignEEATransactionRequest struct {
	Nonce          hexutil.Uint64  `json:"nonce" example:"0x1" swaggertype:"string"`
	To             *common.Address `json:"to,omitempty" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
	Data           hexutil.Bytes   `json:"data,omitempty" example:"0xfeaeee..." swaggertype:"string"`
	ChainID        hexutil.Big     `json:"chainID" validate:"required" example:"0x1 (mainnet)" swaggertype:"string"`
	PrivateFrom    string          `json:"privateFrom" validate:"required,base64,required_with=PrivateFor PrivacyGroupID" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string        `json:"privateFor,omitempty" validate:"omitempty,min=1,unique,dive,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string          `json:"privacyGroupId,omitempty" validate:"omitempty,base64" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

type ECRecoverRequest struct {
	Data      hexutil.Bytes `json:"data" validate:"required" example:"0xfeaeee..." swaggertype:"string"`
	Signature hexutil.Bytes `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
}

type VerifyEth1SignatureRequest struct {
	Data      hexutil.Bytes  `json:"data" validate:"required" example:"0xfeaeee..." swaggertype:"string"`
	Signature hexutil.Bytes  `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
	Address   common.Address `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
}

type VerifyTypedDataRequest struct {
	TypedData SignTypedDataRequest `json:"data" validate:"required" swaggertype:"string"`
	Signature hexutil.Bytes        `json:"signature" validate:"required" example:"0x6019a3c8..." swaggertype:"string"`
	Address   common.Address       `json:"address" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
}

type Eth1AccountResponse struct {
	Address   common.Address    `json:"address" example:"0x664895b5fE3ddf049d2Fb508cfA03923859763C6" swaggertype:"string"`
	Key       KeyResponse       `json:"key"`
	Tags      map[string]string `json:"tags,omitempty"`
	CreatedAt time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt time.Time         `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	Disabled  bool              `json:"disabled" example:"false"`
}
