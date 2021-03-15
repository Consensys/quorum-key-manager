package manifest

import (
	"encoding/json"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
)

// Manifest for a store
type Manifest struct {
	// Kind  of store
	Kind types.Kind `json:"kind"`

	// Version
	Version string `json:"version"`

	// Name of the store
	Name string `json:"name"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// Specs specifications of a store
	Specs json.RawMessage `json:"specs"`
}

func (mnfst *Manifest) UnmarshalSpecs(specs interface{}) error {
	return json.Unmarshal(mnfst.Specs, specs)
}
