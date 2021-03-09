package entities

import "time"

// CryptoOperation type of crypto operation
type CryptoOperation string

const (
	Signing    = "signing"
	Encryption = "encryption"
)

// RecoveryPolicy policies for recovering a deleted item
type RecoveryPolicy string

// Attributes are user set configuration and information attached to stored item
type Attributes struct {
	// Operations supported by a stored item (e.g sign, encrypt...)
	Operations []CryptoOperation

	// Disabled wether item is disabled
	Disabled bool

	// TTL
	TTL time.Duration

	// Recovery policy about a key after being deleted before being destroyed
	Recovery *Recovery

	// Tags attached to a stored item
	Tags map[string]string
}

type Recovery struct {
	// Policy for recovery
	Policy RecoveryPolicy

	// Period for recovery
	Period time.Time
}
