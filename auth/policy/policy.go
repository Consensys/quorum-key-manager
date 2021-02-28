package policy

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/store/types"
)

// Policy is a set of permissions to perform actions

// Each client request has a particular set of policies
type Policy struct {
	// Name of the policy
	Name string

	// Type of policy
	Type PolicyType

	// Endorsement of the policy proper to each type of policy
	Endorsement interface{}
}

type PolicyType uint32

const (
	PolicyTypeUnknown PolicyType = iota
	PolicyTypeJSONRPC
	PolicyTypeStore
)

func (p PolicyType) String() string {
	switch p {
	case PolicyTypeJSONRPC:
		return "jsonrpc"
	case PolicyTypeStore:
		return "rgp"
	default:
		return "unknown"
	}
}

type Policies map[PolicyType][]*Policy

func (policies *Policies) IsStoreAuthorized(info *types.StoreInfo) error {
	for _, policy := range (*policies)[PolicyTypeStore] {
		if err := policy.Endorsement.(*StoreEndorsement).IsAuthorized(info); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}

func (policies *Policies) IsJSONRPCAuthorized(method string) error {
	for _, plcy := range (*policies)[PolicyTypeJSONRPC] {
		if err := plcy.Endorsement.(*JSONRPCEndorsement).IsAuthorized(method); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}
