package manager

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/auth/types"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
)

var GroupKind manifest.Kind = "Group"

type GroupSpecs struct {
	Policies []string `json:"policies"`
}

var PolicyKind manifest.Kind = "Policy"

type PolicySpecs struct {
	Statements []*types.Statement `json:"statements"`
}
