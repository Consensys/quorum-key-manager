package manifest

import (
	"encoding/json"

	json2 "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

type Kind string

// Manifest for a store
type Manifest struct {
	// Kind  of store
	Kind Kind `json:"kind"`

	// Version
	Version string `json:"version"`

	// Name of the store
	Name string `json:"name"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// Specs specifications of a store
	Specs interface{} `json:"specs"`
}

func (mnfst *Manifest) UnmarshalSpecs(specs interface{}) error {
	bdata, err := json.Marshal(json2.RecursiveToJSON(mnfst.Specs))
	if err != nil {
		return err
	}
	
	return json.Unmarshal(bdata, specs)
}
