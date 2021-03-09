package errors

import (
	"fmt"
)

const (
	// Connection Errors (class 08XXX)
	Connection               uint64 = 8 << 12
	HashicorpVaultConnection        = Connection + 8<<8 // Service Connection error (subclass 088XX)

	// Data Errors (class 42XXX)
	Data             uint64 = 4<<16 + 2<<12
	Encoding                = Data + 1<<8 // Invalid Encoding (subclass 421XX)
	InvalidFormat           = Data + 3<<8 // Invalid format (subclass 423XX)
	InvalidParameter        = Data + 4<<8 // Invalid parameter provided (subclass 424XX)

	// Storage Error (class DBXXX)
	Storage            uint64 = 13<<16 + 11<<12
	NotFound                  = Storage + 2<<8         // Not found (subclass DB2XX)
	ConstraintViolated        = Storage + 1<<8         // Storage constraint violated (subclass DB1XX)
	AlreadyExists             = ConstraintViolated + 1 // A resource with same index already exists (code DB101)

	// Configuration errors (class F0XXX)
	Config uint64 = 15 << 16

	// Internal errors (class FFXXX)
	Internal uint64 = 15<<16 + 15<<12
)

var NotImplementedError = fmt.Errorf("not implemented")

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

// HashicorpVaultConnectionError is raised when failing to perform on Hashicorp Vault
func HashicorpVaultConnectionError(format string, a ...interface{}) *Error {
	return Errorf(HashicorpVaultConnection, format, a...)
}

// IsHashicorpVaultConnectionError indicate whether an error is a Hashicorp Vault connection error
func IsHashicorpVaultConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), HashicorpVaultConnection)
}

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
}

// NoDataFoundError is raised when accessing a missing Data
func NotFoundError(format string, a ...interface{}) *Error {
	return Errorf(NotFound, format, a...)
}

// IsNotFoundError indicate whether an error is a no Data found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotFound)
}

// InvalidFormatError is raised when a Data does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *Error {
	return Errorf(InvalidFormat, format, a...)
}

// IsInvalidFormatError indicate whether an error is an invalid format error
func IsInvalidFormatError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidFormat)
}

// InvalidParameterError is raised when a provided parameter invalid
func InvalidParameterError(format string, a ...interface{}) *Error {
	return Errorf(InvalidParameter, format, a...)
}

// IsInvalidParameterError indicate whether an error is an invalid parameter error
func IsInvalidParameterError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidParameter)
}

// AlreadyExistsError is raised when a Data constraint has been violated
func AlreadyExistsError(format string, a ...interface{}) *Error {
	return Errorf(AlreadyExists, format, a...)
}

// IsAlreadyExistsError indicate whether an error is an already exists error
func IsAlreadyExistsError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AlreadyExists)
}

// ConfigError is raised when an error is encountered while loading configuration
func ConfigError(format string, a ...interface{}) *Error {
	return Errorf(Config, format, a...)
}
