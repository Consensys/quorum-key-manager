package auth

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

type Policies map[PolicyType][]*Policy

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
		return "store"
	default:
		return "unknown"
	}
}
