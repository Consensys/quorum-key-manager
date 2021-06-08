package types

import (
	manifest2 "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
)

const (
	HashicorpSecrets manifest2.Kind = "HashicorpSecrets"
	AKVSecrets       manifest2.Kind = "AKVSecrets"
	AWSSecrets       manifest2.Kind = "AWSSecrets"
	KMSSecrets       manifest2.Kind = "KMSSecrets"
)
