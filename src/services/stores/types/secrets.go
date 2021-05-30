package types

import (
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
)

const (
	HashicorpSecrets manifest.Kind = "HashicorpSecrets"
	AKVSecrets       manifest.Kind = "AKVSecrets"
	AWSSecrets       manifest.Kind = "AWSSecrets"
	KMSSecrets       manifest.Kind = "KMSSecrets"
)
