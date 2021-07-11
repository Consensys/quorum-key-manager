package policy

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

var GroupKind manifest.Kind = "Group"

type GroupSpecs struct {
	Policies []string `json:"policies"`
}

var PolicyKind manifest.Kind = "Policy"

type PolicySpecs struct {
	Statements []*types.Statement `json:"statements"`
}
