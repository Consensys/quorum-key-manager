package types

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
)

const (
	HashicorpSecrets manifest.Kind = "HashicorpSecrets"
	AKVSecrets       manifest.Kind = "AKVSecrets"
	AWSSecrets       manifest.Kind = "AWSSecrets"
	KMSSecrets       manifest.Kind = "KMSSecrets"
)
