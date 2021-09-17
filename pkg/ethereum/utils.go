package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/signer/core"
)

const (
	EIP712DomainLabel = "EIP712Domain"
)

func GetEIP712EncodedData(typedData *core.TypedData) ([]byte, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	domainSeparatorHash, err := typedData.HashStruct(EIP712DomainLabel, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash)), nil
}

func GetEIP191EncodedData(msg []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%v", len(msg), string(msg)))
}
