package eth

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"

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

func getEIP191EncodedData(msg []byte) string {
	return fmt.Sprintf("\x19Ethereum Signed Message\n%d%v", len(msg), hexutil.Encode(msg))
}

func newEthAccount(key *entities.Key, attr *entities.Attributes) *entities.ETHAccount {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &entities.ETHAccount{
		KeyID:               key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey),
		Tags:                attr.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
		Metadata: &entities.Metadata{
			Disabled:  key.Metadata.Disabled,
			CreatedAt: key.Metadata.CreatedAt,
			UpdatedAt: key.Metadata.UpdatedAt,
		},
	}
}
