package errors

const (
	// Data Errors (class 42XXX)
	Data             uint64 = 4<<16 + 2<<12
	Encoding                = Data + 1<<8 // Invalid Encoding (subclass 421XX)
	InvalidFormat           = Data + 3<<8 // Invalid format (subclass 423XX)
	InvalidParameter        = Data + 4<<8 // Invalid parameter provided (subclass 424XX)
)

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
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

// InvalidRequestError is raised when a client request is invalid
func InvalidRequestError(format string, a ...interface{}) *Error {
	return Errorf(InvalidRequest, format, a...)
}

// IsInvalidRequestError indicate whether an error is an invalid request error
func IsInvalidRequestError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidRequest)
}
