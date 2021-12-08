package entities

const (
	RoleKind  string = "Role"
	NodeKind  string = "Node"
	StoreKind string = "Store"
	VaultKind string = "Vault"
)

type Manifest struct {
	Kind           string      `yaml:"kind" validate:"required,isManifestKind" example:"store"`
	Name           string      `yaml:"name" validate:"required" example:"my-store"`
	ResourceType   string      `yaml:"type,omitempty" example:"ethereum"`
	Specs          interface{} `yaml:"specs" validate:"required"`
	AllowedTenants []string    `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
}
