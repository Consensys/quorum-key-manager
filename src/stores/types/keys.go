package types

import (
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

const (
	HashicorpKeys manifest.Kind = "HashicorpKeys"
	AKVKeys       manifest.Kind = "AKVKeys"
	AWSKeys       manifest.Kind = "AWSKeys"
	LocalKeys     manifest.Kind = "LocalKeys"
)
