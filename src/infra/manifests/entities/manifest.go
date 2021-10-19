package manifest

//TODO: Split this file into the different domains where the types belong

type Kind string
type StoreType string
type VaultType string

const (
	Role  Kind = "Role"
	Node  Kind = "Node"
	Store Kind = "Store"

	Ethereum StoreType = "Ethereum"
	Keys     StoreType = "Keys"
	Secrets  StoreType = "Secrets"

	LocalEthereum VaultType = "Ethereum"
	HashicorpKeys VaultType = "HashicorpKeys"
	AKVKeys       VaultType = "AKVKeys"
	AWSKeys       VaultType = "AWSKeys"
	LocalKeys     VaultType = "LocalKeys"

	HashicorpSecrets VaultType = "HashicorpSecrets"
	AKVSecrets       VaultType = "AKVSecrets"
	AWSSecrets       VaultType = "AWSSecrets"
)

// Manifest for a store
type Manifest struct {
	// Kind of item (Store, Node,...)
	Kind Kind `json:"kind" validate:"required"`

	// Name of the item
	Name string `json:"name" validate:"required"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// AllowedTenants are the tenants allowed to access this item. Public if empty
	AllowedTenants []string `json:"allowedTenants" yaml:"allowedTenants"`

	// Specs specifications of a store
	Specs interface{} `json:"specs" validate:"required"`
}
