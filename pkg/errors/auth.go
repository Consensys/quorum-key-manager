package errors

const (
	// Authentication Errors (class 09XXX)
	InvalidAuthentication uint64 = 9 << 12
	Unauthorized                 = InvalidAuthentication + 1 // Invalid request credentials (code 09001)
)

// UnauthorizedError is raised when authentication credentials are invalid
func UnauthorizedError(format string, a ...interface{}) *Error {
	return Errorf(Unauthorized, format, a...)
}

func IsUnauthorizedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Unauthorized)
}
