package formatters

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FormatEth1AccResponse(key *entities.ETH1Account) *types.Eth1Response {
	return &types.Eth1Response{
		ID:                  key.ID,
		Address:             key.Address,
		PublicKey:           hexutil.Encode(key.PublicKey),
		CompressedPublicKey: hexutil.Encode(key.CompressedPublicKey),
		Disabled:            key.Metadata.Disabled,
		CreatedAt:           key.Metadata.CreatedAt,
		UpdatedAt:           key.Metadata.UpdatedAt,
		ExpireAt:            key.Metadata.ExpireAt,
		DeletedAt:           key.Metadata.DeletedAt,
		DestroyedAt:         key.Metadata.DestroyedAt,
	}
}
