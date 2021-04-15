package hashicorp

import (
	"encoding/json"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/hashicorp/vault/api"
)

func parseResponse(hashicorpSecret *api.Secret) *entities.Key {
	key := &entities.Key{
		ID:        hashicorpSecret.Data[idLabel].(string),
		PublicKey: hashicorpSecret.Data[publicKeyLabel].(string),
		Algo: &entities.Algorithm{
			Type:          hashicorpSecret.Data[algorithmLabel].(string),
			EllipticCurve: hashicorpSecret.Data[curveLabel].(string),
		},
		Metadata: &entities.Metadata{
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

	return key
}
