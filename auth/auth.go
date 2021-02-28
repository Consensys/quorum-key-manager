package auth

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/auth/policy"
)

// Auth is a piece a contextual auth information attached to every client request after
// authentication succeeded (e.g. JWT based auth, Mutual TLS auth, etc.)

// It shoud be generates by an authentication middleware responsible to validate authentication
// and then attach Auth to the request context so it can flow to sub-sequent components
// that can perform authorization/permissions checks (e.g. JSON-RPC methods checks, store access check, etc.)
type Auth struct {
	// ID of the authenticated client
	ID string

	// Policies associated with the authenticated client
	Policies map[policy.PolicyType][]*policy.Policy

	// Metadata is arbitrary string-type metadata
	Metadata map[string]string

	// Raw is JSON-encodable data that is stored with the auth struct (e.g. JWT Token)
	Raw map[string]interface{}
}

func (auth *Auth) IsStoreAuthorized(storeName string) error {
	for _, plcy := range auth.Policies[policy.PolicyTypeStore] {
		if err := plcy.Endorsement.(*policy.StoreEndorsement).IsAuthorized(storeName); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}

func (auth *Auth) IsJSONRPCAuthorized(method string) error {
	for _, plcy := range auth.Policies[policy.PolicyTypeJSONRPC] {
		if err := plcy.Endorsement.(*policy.JSONRPCEndorsement).IsAuthorized(method); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}
