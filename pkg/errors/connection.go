package errors

const (
	// Connection Errors (class 01XXX)
	Connection     uint64 = 1 << 12
	AKV                   = Connection + 1<<8 // AKV Connection error (subclass 011XX)
	HashicorpVault        = Connection + 2<<8 // Hashicorp Connection error (subclass 012XX)
	AWS                   = Connection + 3<<8 // AWS Connection error (subclass 013XX)

	// Invalid Request Errors (class 02XXX)
	InvalidRequest   uint64 = 2 << 12
	Unauthorized            = InvalidRequest + 1<<8 // Unauthorized error (subclass 021XX)
	NotSupported            = InvalidRequest + 2<<8 // NotSupported error (subclass 022XX)
	NotImplemented          = InvalidRequest + 3<<8 // NotImplemented error (subclass 023XX)
	InvalidFormat           = InvalidRequest + 4<<8 // Invalid format (subclass 024XX)
	InvalidParameter        = InvalidRequest + 5<<8 // Invalid parameter provided (subclass 025XX)
)

// HashicorpVaultError is raised when failing to perform on Hashicorp Vault
func HashicorpVaultError(format string, a ...interface{}) *Error {
	return Errorf(HashicorpVault, format, a...)
}

// IsHashicorpVaultError indicate whether an error is a Hashicorp Vault connection error
func IsHashicorpVaultError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), HashicorpVault)
}

// AKVError is raised when failing to perform on AKV client
func AKVError(format string, a ...interface{}) *Error {
	return Errorf(AKV, format, a...)
}

// IsAKVError indicate whether an error is a AKV client connection error
func IsAKVError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AKV)
}

// AWSError is raised when failing to perform on AWS client
func AWSError(format string, a ...interface{}) *Error {
	return Errorf(AWS, format, a...)
}

// IsAWSError indicate whether an error is a AWS client connection error
func IsAWSError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AWS)
}

// UnauthorizedError is raised when authentication credentials are invalid
func UnauthorizedError(format string, a ...interface{}) *Error {
	return Errorf(Unauthorized, format, a...)
}

func IsUnauthorizedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Unauthorized)
}

// UnauthorizedError is raised when authentication credentials are invalid
func NotSupportedError(format string, a ...interface{}) *Error {
	return Errorf(NotSupported, format, a...)
}

func IsNotSupportedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotSupported)
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
