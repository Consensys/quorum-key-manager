package manifest

import (
	"encoding/json"

	json2 "github.com/consensys/quorum-key-manager/pkg/json"
)

//TODO: Split this file into the different domains where the types belong

type Kind string

const (
	Role Kind = "Role"

	Node Kind = "Node"

	Ethereum Kind = "Ethereum"

	HashicorpKeys Kind = "HashicorpKeys"
	AKVKeys       Kind = "AKVKeys"
	AWSKeys       Kind = "AWSKeys"
	LocalKeys     Kind = "LocalKeys"

	HashicorpSecrets Kind = "HashicorpSecrets"
	AKVSecrets       Kind = "AKVSecrets"
	AWSSecrets       Kind = "AWSSecrets"
)

// Manifest for a store
type Manifest struct {
	// Kind of item (Store, Node,...)
	Kind Kind `json:"kind" validate:"required"`

	// Version of item
	Version string `json:"version"`

	// Name of the item
	Name string `json:"name" validate:"required"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// AllowedTenants are the tenants allowed to access this item. Public if empty
	AllowedTenants []string `json:"allowedTenants" yaml:"allowedTenants"`

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
