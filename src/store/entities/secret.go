package entities

// Secret
type Secret struct {
	Value    string
	Disabled bool
	Recovery *Recovery
	Metadata *Metadata
	Tags     map[string]string
}
