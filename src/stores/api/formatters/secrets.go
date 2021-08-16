package formatters

import (
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func FormatSecretResponse(secret *entities.Secret) *types.SecretResponse {
	return &types.SecretResponse{
		ID:          secret.ID,
		Version:     secret.Metadata.Version,
		Value:       secret.Value,
		Tags:        secret.Tags,
		Disabled:    secret.Metadata.Disabled,
		CreatedAt:   secret.Metadata.CreatedAt,
		UpdatedAt:   secret.Metadata.UpdatedAt,
		ExpireAt:    secret.Metadata.ExpireAt,
		DeletedAt:   secret.Metadata.DeletedAt,
		DestroyedAt: secret.Metadata.DestroyedAt,
	}
}
