package errors

const (
	// Invalid Request Errors (class 03XXX)
	Data     uint64 = 3 << 12
	Encoding        = Data + 1<<8 // Invalid Encoding (subclass 031XX)
)

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
}
