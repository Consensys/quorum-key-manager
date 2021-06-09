package akv

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

func parseSecretBundle(secretBundle *keyvault.SecretBundle) *entities.Secret {
	secret := &entities.Secret{
		Tags:     common.Tomapstr(secretBundle.Tags),
		Metadata: &entities.Metadata{},
	}
	if secretBundle.Value != nil {
		secret.Value = *secretBundle.Value
	}

	if secretBundle.ID != nil {
		// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/secrets/<secretName>/<secretVersion>
		chunks := strings.Split(*secretBundle.ID, "/")
		secret.Metadata.Version = chunks[len(chunks)-1]
		secret.ID = chunks[len(chunks)-2]
	}
	if expires := secretBundle.Attributes.Expires; expires != nil {
		secret.Metadata.ExpireAt = time.Unix(0, expires.Duration().Nanoseconds()).In(time.UTC)
	}
	if created := secretBundle.Attributes.Created; created != nil {
		secret.Metadata.CreatedAt = time.Unix(0, created.Duration().Nanoseconds()).In(time.UTC)
	}
	if updated := secretBundle.Attributes.Updated; updated != nil {
		secret.Metadata.UpdatedAt = time.Unix(0, updated.Duration().Nanoseconds()).In(time.UTC)
	}
	if enabled := secretBundle.Attributes.Enabled; enabled != nil {
		secret.Metadata.Disabled = !*enabled
	}

	return secret
}
