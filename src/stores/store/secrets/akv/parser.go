package akv

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

func parseDeletedSecretBundle(secretBundle *keyvault.DeletedSecretBundle) *entities.Secret {
	return buildNewSecret(secretBundle.ID, secretBundle.Value, secretBundle.Tags, secretBundle.Attributes)
}

func parseSecretBundle(secretBundle *keyvault.SecretBundle) *entities.Secret {
	return buildNewSecret(secretBundle.ID, secretBundle.Value, secretBundle.Tags, secretBundle.Attributes)
}

func buildNewSecret(id, value *string, tags map[string]*string, attributes *keyvault.SecretAttributes) *entities.Secret {
	secret := &entities.Secret{
		Tags:     common.Tomapstr(tags),
		Metadata: &entities.Metadata{},
	}
	if value != nil {
		secret.Value = *value
	}

	if id != nil {
		// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/secrets/<secretName>/<secretVersion>
		chunks := strings.Split(*id, "/")
		secret.Metadata.Version = chunks[len(chunks)-1]
		secret.ID = chunks[len(chunks)-2]
	}

	if expires := attributes.Expires; expires != nil {
		secret.Metadata.ExpireAt = time.Unix(0, expires.Duration().Nanoseconds()).In(time.UTC)
	}
	if created := attributes.Created; created != nil {
		secret.Metadata.CreatedAt = time.Unix(0, created.Duration().Nanoseconds()).In(time.UTC)
	}
	if updated := attributes.Updated; updated != nil {
		secret.Metadata.UpdatedAt = time.Unix(0, updated.Duration().Nanoseconds()).In(time.UTC)
	}
	if enabled := attributes.Enabled; enabled != nil {
		secret.Metadata.Disabled = !*enabled
	}

	return secret
}
