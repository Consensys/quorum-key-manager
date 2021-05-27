package errors

const (
	// Connection Errors (class 08XXX)
	Connection               uint64 = 8 << 12
	AKVConnection                   = Connection + 7<<8 // Service Connection error (subclass 087XX)
	HashicorpVaultConnection        = Connection + 8<<8 // Service Connection error (subclass 088XX)
	AWSConnection                   = Connection + 9<<8 // Service Connection error (subclass 089XX)

	// Invalid Request Errors (class 09XXX)
	InvalidRequest uint64 = 9 << 12
	Unauthorized          = InvalidRequest + 1    // Invalid request credentials (code 09001)
	NotSupported          = InvalidRequest + 7<<8 // Not supported request (code 097XX)
)

// HashicorpVaultConnectionError is raised when failing to perform on Hashicorp Vault
func HashicorpVaultConnectionError(format string, a ...interface{}) *Error {
	return Errorf(HashicorpVaultConnection, format, a...)
}

// IsHashicorpVaultConnectionError indicate whether an error is a Hashicorp Vault connection error
func IsHashicorpVaultConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), HashicorpVaultConnection)
}

// AKVConnectionError is raised when failing to perform on AKV client
func AKVConnectionError(format string, a ...interface{}) *Error {
	return Errorf(AKVConnection, format, a...)
}

// IsAKVConnectionError indicate whether an error is a AKV client connection error
func IsAKVConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AKVConnection)
}

// AWSConnectionError is raised when failing to perform on AWS client
func AWSConnectionError(format string, a ...interface{}) *Error {
	return Errorf(AWSConnection, format, a...)
}

// IsAWSConnectionError indicate whether an error is a AWS client connection error
func IsAWSConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AWSConnection)
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
