package eth1

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/ethereum/go-ethereum/signer/core"
)

func getEIP712EncodedData(typedData *core.TypedData) (string, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", err
	}

	domainSeparatorHash, err := typedData.HashStruct(formatters.EIP712DomainLabel, typedData.Domain.Map())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash), nil
}

func getEIP191EncodedData(msg string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message\n%d%v", len(msg), msg)
}
