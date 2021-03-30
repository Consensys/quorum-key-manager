package hashicorp

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/hashicorp/vault/api"
)

type hashicorpKey struct {
	ID        string            `json:"id"`
	Curve     string            `json:"curve"`
	Algorithm string            `json:"algorithm"`
	PublicKey string            `json:"publicKey"`
	Namespace string            `json:"namespace,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
}

func parseResponse(data map[string]interface{}) (*entities.Key, error) {
	key := &hashicorpKey{}

	jsonbody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonbody, &key)
	if err != nil {
		return nil, err
	}

	return &entities.Key{
		ID:        key.ID,
		PublicKey: key.PublicKey,
		Algo: &entities.Algorithm{
			Type:          key.Algorithm,
			EllipticCurve: key.Curve,
		},
		// TODO: Add metadata when this is added to the plugin
		Metadata: &entities.Metadata{
			Version:  1,
			Disabled: false,
		},
		Tags: key.Tags,
	}, nil
}

func parseErrorResponse(err error) error {
	httpError, _ := err.(*api.ResponseError)

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
