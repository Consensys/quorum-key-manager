package errors

const (
	Data            = "DA000"
	Encoding        = "DA100"
	CryptoOperation = "DA200"
)

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
}

// CryptoOperationError are raised when failing to perform a crypto operation
func CryptoOperationError(format string, a ...interface{}) *Error {
	return Errorf(CryptoOperation, format, a...)
}

// IsCryptoOperationError indicate whether an error is a CryptoOperationError error
func IsCryptoOperationError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), CryptoOperation)
}
