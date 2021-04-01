package hashicorp

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
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

func parseErrorResponse(httpError *api.ResponseError) error {
	switch httpError.StatusCode {
	case http.StatusNotFound:
		return errors.NotFoundError(httpError.Error())
	case http.StatusBadRequest:
		return errors.InvalidFormatError(httpError.Error())
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(httpError.Error())
	case http.StatusConflict:
		return errors.AlreadyExistsError(httpError.Error())
	default:
		return errors.HashicorpVaultConnectionError(httpError.Error())
	}
}
