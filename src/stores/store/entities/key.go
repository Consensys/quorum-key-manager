package entities

// Key public part of a key
type Key struct {
	ID          string
	PublicKey   []byte
	Algo        *Algorithm
	Metadata    *Metadata
	Tags        map[string]string
	Annotations map[string]string
}
