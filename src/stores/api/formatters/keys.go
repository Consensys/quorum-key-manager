package formatters

import (
	"encoding/base64"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

func FormatKeyResponse(key *entities2.Key) *types2.KeyResponse {
	return &types2.KeyResponse{
		ID:               key.ID,
		PublicKey:        base64.URLEncoding.EncodeToString(key.PublicKey),
		Curve:            string(key.Algo.EllipticCurve),
		SigningAlgorithm: string(key.Algo.Type),
		Tags:             key.Tags,
		Disabled:         key.Metadata.Disabled,
		CreatedAt:        key.Metadata.CreatedAt,
		UpdatedAt:        key.Metadata.UpdatedAt,
		ExpireAt:         key.Metadata.ExpireAt,
		DeletedAt:        key.Metadata.DeletedAt,
		DestroyedAt:      key.Metadata.DestroyedAt,
	}
}
