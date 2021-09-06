package entities

import "github.com/ethereum/go-ethereum/common"

type ETHAccount struct {
	Address             common.Address
	KeyID               string
	PublicKey           []byte
	CompressedPublicKey []byte
	Metadata            *Metadata
	Tags                map[string]string
}
