package errors

import "strings"

const (
	Internal          = "IN000"
	Config            = "IN100"
	DependencyFailure = "IN200"
)

//nolint
var ErrNotImplemented = NotImplementedError("this operation is not yet implemented. Please contact your administrator")
var ErrNotSupported = NotSupportedError("this operation is not supported. Please contact your administrator")

func isErrorClass(code, base string) bool {
	// Delete tailing 0's of the base error code
	for base[len(base)-1] == '0' {
		base = base[:len(base)-1]
	}

	return strings.HasPrefix(code, base)
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
