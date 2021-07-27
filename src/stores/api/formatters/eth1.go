package formatters

import (
	"math/big"

	common2 "github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

func FormatTransaction(tx *types.SignETHTransactionRequest) *ethtypes.Transaction {
	if tx.To == nil {
		return ethtypes.NewContractCreation(uint64(tx.Nonce), tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
	}
	return ethtypes.NewTransaction(uint64(tx.Nonce), *tx.To, tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
}

func FormatPrivateTransaction(tx *types.SignQuorumPrivateTransactionRequest) *quorumtypes.Transaction {
	if tx.To == nil {
		return quorumtypes.NewContractCreation(uint64(tx.Nonce), tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
	}
	return quorumtypes.NewTransaction(uint64(tx.Nonce), *tx.To, tx.Value.ToInt(), uint64(tx.GasLimit), tx.GasPrice.ToInt(), tx.Data)
}

func FormatEEATransaction(tx *types.SignEEATransactionRequest) (*ethtypes.Transaction, *ethereum.PrivateArgs) {
	privateArgs := &ethereum.PrivateArgs{
		PrivateFrom:    &tx.PrivateFrom,
		PrivateFor:     &tx.PrivateFor,
		PrivateType:    common2.ToPtr(PrivateTxTypeRestricted).(*string),
		PrivacyGroupID: &tx.PrivacyGroupID,
	}

	if tx.To == nil {
		return ethtypes.NewContractCreation(uint64(tx.Nonce), big.NewInt(0), uint64(0), big.NewInt(0), tx.Data), privateArgs
	}
	return ethtypes.NewTransaction(uint64(tx.Nonce), *tx.To, big.NewInt(0), uint64(0), big.NewInt(0), tx.Data), privateArgs
}

func FormatEth1AccResponse(eth1Acc *entities.ETH1Account) *types.Eth1AccountResponse {
	return &types.Eth1AccountResponse{
		Address:   eth1Acc.Address,
		Key:       *FormatKeyResponse(eth1Acc.Key),
		Tags:      eth1Acc.Tags,
		Disabled:  eth1Acc.Metadata.Disabled,
		CreatedAt: eth1Acc.Metadata.CreatedAt,
		UpdatedAt: eth1Acc.Metadata.UpdatedAt,
		DeletedAt: eth1Acc.Metadata.DeletedAt,
	}
}
