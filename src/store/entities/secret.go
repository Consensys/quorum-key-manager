package entities

// Secret
type Secret struct {
	Value    string
	Recovery *Recovery
	Metadata *Metadata
	Tags     map[string]string
}
