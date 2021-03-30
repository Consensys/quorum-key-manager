package entities

type Key struct {
	ID        string
	PublicKey string
	Algo      *Algorithm
	Metadata  *Metadata
	Tags      map[string]string
}
