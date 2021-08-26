package formatters

import (
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func FormatSecretResponse(secret *entities.Secret) *types.SecretResponse {
	resp := &types.SecretResponse{
		ID:        secret.ID,
		Version:   secret.Metadata.Version,
		Value:     secret.Value,
		Tags:      secret.Tags,
		Disabled:  secret.Metadata.Disabled,
		CreatedAt: secret.Metadata.CreatedAt,
		UpdatedAt: secret.Metadata.UpdatedAt,
	}

	if !secret.Metadata.DeletedAt.IsZero() {
		resp.DeletedAt = &secret.Metadata.DeletedAt
	}

	return resp
}
