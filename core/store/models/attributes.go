package models

import "time"

// CryptoOperation type of crypto operation
type CryptoOperation int

// RecoveryPolicy policies for recovering a deleted item
type RecoveryPolicy int

// Attributes are user set configuration and information attached to stored item
type Attributes struct {
	// Operations supported by a stored item (e.g sign, verify, encrypt...)
	Operations []CryptoOperation

	// Enabled wether item is enabled
	Enabled bool

	// ExpireAt expiration date
	ExpireAt time.Time

	// Recovery policy about a key after being deleted before being destroyed
	Recovery struct {
		// Policy for recovery
		Policy RecoveryPolicy

		// Period for recovery
		Period time.Time
	}

	// Tags attached to a stored item
	Tags map[string]string
}
