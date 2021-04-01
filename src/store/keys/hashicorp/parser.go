package hashicorp

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/hashicorp/vault/api"
)

func parseResponse(hashicorpSecret *api.Secret) *entities.Key {
	return &entities.Key{
		ID:        hashicorpSecret.Data[idLabel].(string),
		PublicKey: hashicorpSecret.Data[publicKeyLabel].(string),
		Algo: &entities.Algorithm{
			Type:          hashicorpSecret.Data[algorithmLabel].(string),
			EllipticCurve: hashicorpSecret.Data[curveLabel].(string),
		},
		// TODO: Add metadata when this is added to the plugin
		Metadata: &entities.Metadata{
			Version:  1,
			Disabled: false,
		},
		Tags: hashicorpSecret.Data[tagsLabel].(map[string]string),
	}
}
