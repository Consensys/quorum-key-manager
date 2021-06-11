package errors

const (
	Internal string = "IN"

	Config            = Internal + "1"
	DependencyFailure = Internal + "2"
)

//nolint
var ErrNotImplemented = NotImplementedError("this operation is not yet implemented. Please contact your administrator")
var ErrNotSupported = NotSupportedError("this operation is not supported. Please contact your administrator")

func isErrorClass(code, base string) bool {
	return code == base
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
