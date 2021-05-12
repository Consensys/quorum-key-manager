package entities

type ETH1Account struct {
	ID                  string
	Address             string
	Metadata            *Metadata
	PublicKey           string
	CompressedPublicKey string
	Tags                map[string]string
}
