package auth

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/auth/policy"
)

var SystemID = "system"

// SystemAuth is a default Auth for system action
var SystemAuth = &Auth{
	ID: SystemID,
	Policies: map[policy.PolicyType][]*policy.Policy{
		policy.PolicyTypeJSONRPC: []*policy.Policy{
			SystemJSONRPCPolicy,
		},
		policy.PolicyTypeStore: []*policy.Policy{
			SystemStorePolicy,
		},
	},
}

var SystemJSONRPCPolicy = &policy.Policy{
	Name: fmt.Sprintf("%v.jsonrpc", SystemID),
	Type: policy.PolicyTypeJSONRPC,
	Endorsement: &policy.JSONRPCEndorsement{
		Scopes: []*policy.JSONRPCScope{
			&policy.JSONRPCScope{
				Service: "*",
				Method:  "*",
			},
		},
	},
}

var SystemStorePolicy = &policy.Policy{
	Name: fmt.Sprintf("%v.store", SystemID),
	Type: policy.PolicyTypeJSONRPC,
	Endorsement: &policy.StoreEndorsement{
		Names: []string{"*"},
	},
}
