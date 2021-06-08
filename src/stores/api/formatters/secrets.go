package formatters

import (
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

func FormatSecretResponse(secret *entities2.Secret) *types2.SecretResponse {
	return &types2.SecretResponse{
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
