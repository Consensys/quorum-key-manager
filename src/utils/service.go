package utils

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

type Utilities interface {
	// Verify verifies the signature belongs to the corresponding key
	Verify(pubKey, data, sig []byte, algo *entities.Algorithm) error

	// ECRecover returns the Ethereum address from a signature and data
	ECRecover(data, sig []byte) (common.Address, error)

	// VerifyMessage verifies that a message signature belongs to a given address
	VerifyMessage(addr common.Address, data, sig []byte) error

	// VerifyTypedData verifies that a typed data signature belongs to a given address
	VerifyTypedData(addr common.Address, typedData *core.TypedData, sig []byte) error
}
