package errors

const (
	Data     = "DA000"
	Encoding = "DA100"
)

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
}
