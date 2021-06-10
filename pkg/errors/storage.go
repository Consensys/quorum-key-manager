package errors

const (
	// Storage Error (class DBXXX)
	Storage        uint64 = 13<<16 + 11<<12
	NotFound              = Storage + 1<<8 // Not found (subclass DB1XX)
	AlreadyExists         = Storage + 2<<8 // A resource with same index already exists (code DB201)
	StatusConflict        = Storage + 3<<8 // A resource exists with invalid status (code DB301)
)

// NoDataFoundError is raised when accessing a missing Data
func NotFoundError(format string, a ...interface{}) *Error {
	return Errorf(NotFound, format, a...)
}

// IsNotFoundError indicate whether an error is a no Data found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotFound)
}

// AlreadyExistsError is raised when a Data constraint has been violated
func AlreadyExistsError(format string, a ...interface{}) *Error {
	return Errorf(AlreadyExists, format, a...)
}

// IsAlreadyExistsError indicate whether an error is an already exists error
func IsAlreadyExistsError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AlreadyExists)
}

// AlreadyExistsError is raised when a Data constraint has been violated
func StatusConflictError(format string, a ...interface{}) *Error {
	return Errorf(StatusConflict, format, a...)
}

// IsAlreadyExistsError indicate whether an error is an already exists error
func IsStatusConflictError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), StatusConflict)
}
