package entities

import "github.com/ethereum/go-ethereum/common"

type ETH1Account struct {
	ID                  string
	Address             common.Address
	Metadata            *Metadata
	PublicKey           []byte
	CompressedPublicKey []byte
	Tags                map[string]string
}
