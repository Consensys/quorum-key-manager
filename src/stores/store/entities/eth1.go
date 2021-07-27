package entities

import "github.com/ethereum/go-ethereum/common"

type ETH1Account struct {
	Address             common.Address
	KeyID               string
	Key                 *Key
	Metadata            *Metadata
	Tags                map[string]string
}
