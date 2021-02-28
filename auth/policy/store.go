package policy

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/store/types"
)

// StoreEndorsement list store accesses
type StoreEndorsement struct {
	Names []string
}

// IsAuthorized indicates wether
func (e *StoreEndorsement) IsAuthorized(info *types.StoreInfo) error {
	// TODO: check if storeName is authorized
	// - Basic implementation can simply check if storeName is in the slice
	// - More advanced implementation can consider some regexp
	return fmt.Errorf("not implemented error")
}
