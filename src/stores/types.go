package stores

import (
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

const (
	Eth1Account manifest.Kind = "Eth1Account"
)

const (
	HashicorpKeys manifest.Kind = "HashicorpKeys"
	AKVKeys       manifest.Kind = "AKVKeys"
	AWSKeys       manifest.Kind = "AWSKeys"
	LocalKeys     manifest.Kind = "LocalKeys"
)

const (
	HashicorpSecrets manifest.Kind = "HashicorpSecrets"
	AKVSecrets       manifest.Kind = "AKVSecrets"
	AWSSecrets       manifest.Kind = "AWSSecrets"
)
