package entities

const (
	RoleKind  string = "role"
	NodeKind  string = "node"
	StoreKind string = "store"
	VaultKind string = "vault"
)

type Manifest struct {
	Kind  string      `yaml:"kind" validate:"required,isManifestKind" example:"store"`
	Specs interface{} `yaml:"specs" validate:"required"`
}
