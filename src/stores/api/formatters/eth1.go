package formatters

import (
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

const (
	PrivateTxTypeRestricted = "restricted"
	EIP712DomainLabel       = "EIP712Domain"
)

func FormatSignTypedDataRequest(request *types2.SignTypedDataRequest) *signer.TypedData {
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

func FormatTransaction(tx *types2.SignETHTransactionRequest) *ethtypes.Transaction {
	if tx.To == "" {
		return ethtypes.NewContractCreation(tx.Nonce, tx.Value.ToInt(), tx.GasLimit, tx.GasPrice.ToInt(), tx.Data)
	}
	return ethtypes.NewTransaction(tx.Nonce, common.HexToAddress(tx.To), tx.Value.ToInt(), tx.GasLimit, tx.GasPrice.ToInt(), tx.Data)
}

func FormatPrivateTransaction(tx *types2.SignQuorumPrivateTransactionRequest) *quorumtypes.Transaction {
	if tx.To == "" {
		return quorumtypes.NewContractCreation(tx.Nonce, tx.Value.ToInt(), tx.GasLimit, tx.GasPrice.ToInt(), tx.Data)
	}
	return quorumtypes.NewTransaction(tx.Nonce, common.HexToAddress(tx.To), tx.Value.ToInt(), tx.GasLimit, tx.GasPrice.ToInt(), tx.Data)
}

func FormatEEATransaction(tx *types2.SignEEATransactionRequest) (*ethtypes.Transaction, *ethereum.PrivateArgs) {
	privateType := PrivateTxTypeRestricted
	privateArgs := &ethereum.PrivateArgs{
		PrivateFrom:    &tx.PrivateFrom,
		PrivateFor:     &tx.PrivateFor,
		PrivateType:    &privateType,
		PrivacyGroupID: &tx.PrivacyGroupID,
	}

	if tx.To == "" {
		return ethtypes.NewContractCreation(tx.Nonce, big.NewInt(0), uint64(0), big.NewInt(0), tx.Data), privateArgs
	}
	return ethtypes.NewTransaction(tx.Nonce, common.HexToAddress(tx.To), big.NewInt(0), uint64(0), big.NewInt(0), tx.Data), privateArgs
}

func FormatEth1AccResponse(key *entities2.ETH1Account) *types2.Eth1AccountResponse {
	return &types2.Eth1AccountResponse{
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
