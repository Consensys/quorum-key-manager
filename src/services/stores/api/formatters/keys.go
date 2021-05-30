package formatters

import (
	"encoding/base64"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
)

func FormatKeyResponse(key *entities.Key) *types.KeyResponse {
	return &types.KeyResponse{
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
