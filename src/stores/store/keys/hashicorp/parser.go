package hashicorp

import (
	"encoding/base64"
	"encoding/json"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"

	"github.com/hashicorp/vault/api"
)

func parseResponse(hashicorpSecret *api.Secret) (*entities2.Key, error) {
	pubKey, err := base64.URLEncoding.DecodeString(hashicorpSecret.Data[publicKeyLabel].(string))
	if err != nil {
		return nil, errors.HashicorpVaultConnectionError("failed to decode public key")
	}

	key := &entities2.Key{
		ID:        hashicorpSecret.Data[idLabel].(string),
		PublicKey: pubKey,
		Algo: &entities2.Algorithm{
			Type:          entities2.KeyType(hashicorpSecret.Data[algorithmLabel].(string)),
			EllipticCurve: entities2.Curve(hashicorpSecret.Data[curveLabel].(string)),
		},
		Metadata: &entities2.Metadata{
			Version:  hashicorpSecret.Data[versionLabel].(json.Number).String(),
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
