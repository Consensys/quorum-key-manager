package formatters

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func FormatKeyResponse(key *entities.Key) *types.KeyResponse {
	return &types.KeyResponse{
		ID:               key.ID,
		PublicKey:        key.PublicKey,
		Curve:            key.Algo.EllipticCurve,
		SigningAlgorithm: key.Algo.Type,
		Tags:             key.Tags,
		Version:          key.Metadata.Version,
		Disabled:         key.Metadata.Disabled,
		CreatedAt:        key.Metadata.CreatedAt,
		UpdatedAt:        key.Metadata.UpdatedAt,
		ExpireAt:         key.Metadata.ExpireAt,
		DeletedAt:        key.Metadata.DeletedAt,
		DestroyedAt:      key.Metadata.DestroyedAt,
	}
}
