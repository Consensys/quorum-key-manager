package formatters

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
)

func FormatSecretResponse(secret *entities.Secret) *types.SecretResponse {
	return &types.SecretResponse{
		ID:          secret.ID,
		Value:       secret.Value,
		Tags:        secret.Tags,
		Version:     secret.Metadata.Version,
		Disabled:    secret.Metadata.Disabled,
		CreatedAt:   secret.Metadata.CreatedAt,
		UpdatedAt:   secret.Metadata.UpdatedAt,
		ExpireAt:    secret.Metadata.ExpireAt,
		DeletedAt:   secret.Metadata.DeletedAt,
		DestroyedAt: secret.Metadata.DestroyedAt,
	}
}
