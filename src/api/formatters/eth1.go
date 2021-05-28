package formatters

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

const (
	PrivateTxTypeRestricted = "restricted"
	EIP712DomainLabel       = "EIP712Domain"
)

func FormatSignTypedDataRequest(request *types.SignTypedDataRequest) *signer.TypedData {
	typedData := &signer.TypedData{
		Types: signer.Types{
			EIP712DomainLabel: []signer.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "version", Type: "string"},
			},
		},
		PrimaryType: request.MessageType,
		Domain: signer.TypedDataDomain{
			Name:              request.DomainSeparator.Name,
			Version:           request.DomainSeparator.Version,
			ChainId:           math.NewHexOrDecimal256(request.DomainSeparator.ChainID),
			VerifyingContract: request.DomainSeparator.VerifyingContract,
			Salt:              request.DomainSeparator.Salt,
		},
		Message: request.Message,
	}

	for typeName, requestTypes := range request.Types {
		var typesDefinition []signer.Type
		for _, typeDefRequest := range requestTypes {
			typesDefinition = append(typesDefinition, signer.Type{
				Name: typeDefRequest.Name,
				Type: typeDefRequest.Type,
			})
		}
		typedData.Types[typeName] = typesDefinition
	}

	if request.DomainSeparator.VerifyingContract != "" {
		typedData.Types[EIP712DomainLabel] = append(typedData.Types[EIP712DomainLabel], signer.Type{Name: "verifyingContract", Type: "address"})
	}

	if request.DomainSeparator.Salt != "" {
		typedData.Types[EIP712DomainLabel] = append(typedData.Types[EIP712DomainLabel], signer.Type{Name: "salt", Type: "string"})
	}

	return typedData
}

func FormatEth1AccResponse(key *entities.ETH1Account) *types.Eth1AccountResponse {
	return &types.Eth1AccountResponse{
		ID:                  key.ID,
		Address:             key.Address,
		PublicKey:           hexutil.Encode(key.PublicKey),
		CompressedPublicKey: hexutil.Encode(key.CompressedPublicKey),
		Tags:                key.Tags,
		Disabled:            key.Metadata.Disabled,
		CreatedAt:           key.Metadata.CreatedAt,
		UpdatedAt:           key.Metadata.UpdatedAt,
		ExpireAt:            key.Metadata.ExpireAt,
		DeletedAt:           key.Metadata.DeletedAt,
		DestroyedAt:         key.Metadata.DestroyedAt,
	}
}
