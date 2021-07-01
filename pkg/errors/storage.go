package errors

const (
	Storage        = "ST000"
	NotFound       = "ST100"
	AlreadyExists  = "ST200"
	StatusConflict = "ST300"
)

// NotFoundError is raised when accessing a missing Data
func NotFoundError(format string, a ...interface{}) *Error {
	return Errorf(NotFound, format, a...)
}

// IsNotFoundError indicate whether an error is a no Data found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotFound)
}

// AlreadyExistsError indicates a resource already exists
func AlreadyExistsError(format string, a ...interface{}) *Error {
	return Errorf(AlreadyExists, format, a...)
}

// IsAlreadyExistsError indicate whether an error is an already exists error
func IsAlreadyExistsError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AlreadyExists)
}

// StatusConflictError indicate a status conflict
func StatusConflictError(format string, a ...interface{}) *Error {
	return Errorf(StatusConflict, format, a...)
}

// IsStatusConflictError indicate whether an error is a status conflict error
func IsStatusConflictError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), StatusConflict)
}
