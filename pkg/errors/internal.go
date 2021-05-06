package errors

import (
	"fmt"
)

const (
	// Configuration errors (class F0XXX)
	Config uint64 = 15 << 16

	// Internal errors (class FFXXX)
	Internal uint64 = 15<<16 + 15<<12
)

//nolint
var ErrNotImplemented = fmt.Errorf("not implemented")
var ErrNotSupported = fmt.Errorf("not supported")

func isErrorClass(code, base uint64) bool {
	// Error codes have a 5 hex representation (<=> 20 bits representation)
	//  - (code^base)&255<<12 compute difference between 2 first nibbles (bits 13 to 20)
	//  - (code^base)&(base&15<<8) compute difference between 3rd nibble in case base 3rd nibble is non zero (bits 9 to 12)
	return (code^base)&(255<<12+15<<8&base) == 0
}

// InternalError is raised when an unknown exception is met
func InternalError(format string, a ...interface{}) *Error {
	return Errorf(Internal, format, a...)
}

// IsInternalError indicate whether an error is an Internal error
func IsInternalError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Internal)
}

// ConfigError is raised when an error is encountered while loading configuration
func ConfigError(format string, a ...interface{}) *Error {
	return Errorf(Config, format, a...)
}

func DependencyFailureError(format string, a ...interface{}) *Error {
	return Errorf(Internal, format, a...)
}

func IsDependencyFailureError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Internal)
}
