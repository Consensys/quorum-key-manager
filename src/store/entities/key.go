package entities

// Key public part of a key
type Key struct {
	ID        string
	PublicKey string
	Algo      *Algorithm
	Metadata  *Metadata
	Tags      map[string]string
}
