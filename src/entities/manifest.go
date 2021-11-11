package entities

const (
	RoleKind  string = "Role"
	NodeKind  string = "Node"
	StoreKind string = "Store"
)

type Manifest struct {
	Kind  string      `yaml:"kind" validate:"required,isManifestKind" example:"store"`
	Specs interface{} `yaml:"specs" validate:"required"`
}
