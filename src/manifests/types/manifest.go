package manifest

import (
	"encoding/json"

	json2 "github.com/consensys/quorum-key-manager/pkg/json"
)

type Kind string

// Manifest for a store
type Manifest struct {
	// Kind  of store
	Kind Kind `json:"kind" validate:"required"`

	// Version
	Version string `json:"version"`

	// Name of the store
	Name string `json:"name" validate:"required"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// Specs specifications of a store
	Specs interface{} `json:"specs" validate:"required"`
}

func (mnfst *Manifest) UnmarshalSpecs(specs interface{}) error {
	return UnmarshalSpecs(mnfst.Specs, specs)
}

func UnmarshalSpecs(src, dest interface{}) error {
	bdata, err := json.Marshal(json2.RecursiveToJSON(src))
	if err != nil {
		return err
	}

	return json.Unmarshal(bdata, dest)
}
