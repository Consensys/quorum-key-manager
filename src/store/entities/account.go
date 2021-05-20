package entities

type ETH1Account struct {
	ID                  string
	Address             string
	Metadata            *Metadata
	PublicKey           []byte
	CompressedPublicKey []byte
	Tags                map[string]string
}
