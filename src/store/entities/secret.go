package entities

// Secret
type Secret struct {
	ID       string
	Value    string
	Metadata *Metadata
	Tags     map[string]string
}
