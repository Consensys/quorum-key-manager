package policy

import (
	"fmt"
)

// StoreEndorsement list store accesses
type StoreEndorsement struct {
	Names []string
}

// IsAuthorized indicates wether
func (e *StoreEndorsement) IsAuthorized(storeName string) error {
	// TODO: check if storeName is authorized
	// - Basic implementation can simply check if storeName is in the slice
	// - More advanced implementation can consider some regexp
	return fmt.Errorf("not implemented error")
}

