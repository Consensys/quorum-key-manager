package entities

const (
	RoleKind  string = "role"
	NodeKind  string = "node"
	StoreKind string = "store"
	VaultKind string = "vault"
)

type Manifest struct {
	Kind         string      `yaml:"kind" validate:"required,isManifestKind" example:"store"`
	Name         string      `json:"name" yaml:"name" validate:"required" example:"my-store"`
	ResourceType string      `json:"type,omitempty" yaml:"type,omitempty" example:"ethereum"`
	Specs        interface{} `yaml:"specs" validate:"required"`
}
