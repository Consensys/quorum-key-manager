package hashicorp

import (
	"encoding/base64"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/hashicorp/vault/api"
)

func parseAPISecretToKey(hashicorpSecret *api.Secret) (*entities.Key, error) {
	pubKey, err := base64.URLEncoding.DecodeString(hashicorpSecret.Data[publicKeyLabel].(string))
	if err != nil {
		return nil, errors.HashicorpVaultError("failed to decode public key")
	}

	key := &entities.Key{
		ID:        hashicorpSecret.Data[idLabel].(string),
		PublicKey: pubKey,
		Algo: &entities.Algorithm{
			Type:          entities.KeyType(hashicorpSecret.Data[algorithmLabel].(string)),
			EllipticCurve: entities.Curve(hashicorpSecret.Data[curveLabel].(string)),
		},
		Metadata: &entities.Metadata{
			Disabled: false,
		},
		Tags: make(map[string]string),
	}

	if hashicorpSecret.Data[tagsLabel] != nil {
		tags := hashicorpSecret.Data[tagsLabel].(map[string]interface{})
		for k, v := range tags {
			key.Tags[k] = v.(string)
		}
	}

	key.Metadata.CreatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[createdAtLabel].(string))
	key.Metadata.UpdatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[updatedAtLabel].(string))

	return key, nil
}
