package errors

const (
	// Internal Errors (class 04XXX)
	Internal uint64 = 4 << 12

	Config            = Internal + 1<<8 // Config error (subclass 041XX)
	DependencyFailure = Internal + 2<<8 // DependencyFailure error (subclass 042XX)
)

//nolint
var ErrNotImplemented = NotImplementedError("this operation is not yet implemented. Please contact your administrator")
var ErrNotSupported = NotSupportedError("this operation is not supported. Please contact your administrator")

func isErrorClass(code, base uint64) bool {
	// Error codes have a 5 hex representation (<=> 20 bits representation)
	//  - (code^base)&255<<12 compute difference between 2 first nibbles (bits 13 to 20)
	//  - (code^base)&(base&15<<8) compute difference between 3rd nibble in case base 3rd nibble is non zero (bits 9 to 12)
	return (code^base)&(255<<12+15<<8&base) == 0
}

func ConfigError(format string, a ...interface{}) *Error {
	return Errorf(Config, format, a...)
}

func DependencyFailureError(format string, a ...interface{}) *Error {
	return Errorf(DependencyFailure, format, a...)
}

func IsDependencyFailureError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), DependencyFailure)
}

func NotImplementedError(format string, a ...interface{}) *Error {
	return Errorf(NotImplemented, format, a...)
}

func IsNotImplementedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotImplemented)
}
