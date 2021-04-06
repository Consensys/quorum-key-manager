package errors

const (
	// Connection Errors (class 08XXX)
	Connection               uint64 = 8 << 12
	AKVConnection                   = Connection + 7<<8 // Service Connection error (subclass 088XX)
	HashicorpVaultConnection        = Connection + 8<<8 // Service Connection error (subclass 088XX)
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
