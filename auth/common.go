package auth

import (
	"fmt"
)

var SystemID = "system"

// SystemAuth is a default Auth for system action
var SystemAuth = &Auth{
	ID: SystemID,
	Policies: map[PolicyType][]*Policy{
		PolicyTypeJSONRPC: []*Policy{
			SystemJSONRPCPolicy,
		},
		PolicyTypeStore: []*Policy{
			SystemStorePolicy,
		},
	},
}

var SystemJSONRPCPolicy = &Policy{
	Name: fmt.Sprintf("%v.jsonrpc", SystemID),
	Type: PolicyTypeJSONRPC,
	Endorsement: &JSONRPCEndorsement{
		Scopes: []*JSONRPCScope{
			&JSONRPCScope{
				Service: "*",
				Method:  "*",
			},
		},
	},
}

var SystemStorePolicy = &Policy{
	Name: fmt.Sprintf("%v.store", SystemID),
	Type: PolicyTypeJSONRPC,
	Endorsement: &StoreEndorsement{
		Names: []string{"*"},
	},
}

// NotAuthenticatedAuth is a default Auth which has no policy attached to it
var NotAuthenticatedAuth = &Auth{ID: "not-authenticated"}
