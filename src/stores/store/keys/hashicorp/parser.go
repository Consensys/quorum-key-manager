package hashicorp

import (
	"encoding/base64"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/hashicorp/vault/api"
)

func parseResponse(hashicorpSecret *api.Secret) (*models.Key, error) {
	pubKey, err := base64.URLEncoding.DecodeString(hashicorpSecret.Data[publicKeyLabel].(string))
	if err != nil {
		return nil, errors.HashicorpVaultError("failed to decode public key")
	}

	key := &models.Key{
		ID:               hashicorpSecret.Data[idLabel].(string),
		PublicKey:        pubKey,
		SigningAlgorithm: hashicorpSecret.Data[algorithmLabel].(string),
		EllipticCurve:    hashicorpSecret.Data[curveLabel].(string),
		Disabled:         false,
		Tags:             make(map[string]string),
	}

	if hashicorpSecret.Data[tagsLabel] != nil {
		tags := hashicorpSecret.Data[tagsLabel].(map[string]interface{})
		for k, v := range tags {
			key.Tags[k] = v.(string)
		}
	}

	key.CreatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[createdAtLabel].(string))
	key.UpdatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[updatedAtLabel].(string))

	return key, nil
}
