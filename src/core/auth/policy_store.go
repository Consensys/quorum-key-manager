package auth

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func NewStorePolicy(name string, storeNames []string) *Policy {
	return &Policy{
		Name: name,
		Type: PolicyTypeStore,
		Endorsement: &StoreEndorsement{
			Names: storeNames,
		},
	}
}

// StoreEndorsement list store accesses
type StoreEndorsement struct {
	Names []string
}

// IsAuthorized indicates wether
func (e *StoreEndorsement) IsAuthorized(storeInfo *entities.StoreInfo) error {
	// TODO: check if storeName is authorized
	// - Basic implementation can simply check if storeName is in the slice
	// - More advanced implementation can consider some regexp
	return fmt.Errorf("not implemented error")
}

// IsStoreAuthorized
func (policies *Policies) IsStoreAuthorized(storeInfo *entities.StoreInfo) error {
	for _, policy := range (*policies)[PolicyTypeStore] {
		if err := policy.Endorsement.(*StoreEndorsement).IsAuthorized(storeInfo); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}
