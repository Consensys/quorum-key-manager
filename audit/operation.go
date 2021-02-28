package audit

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/auth"
)

// Operation
type Operation struct {
	// ID of the operation (automatically generated)
	ID string

	// Type of operation (e.g. sign-eea, sign-secp256k1, get-secret, etc.)
	Type string

	// Time when operation start
	StartTime time.Time

	// Time when operation end
	EndTime time.Time

	// Auth that triggered the operation
	Auth *auth.Auth

	// Data information relative to the operation (e.g. store on which operation is executed, inputs, outputs, etc.)
	Data map[string]interface{}

	// Error in the operation
	Error error

	// Parent operation
	Parent *Operation
}
