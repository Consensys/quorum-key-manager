package errors

const (
	// Storage Error (class DBXXX)
	Storage       uint64 = 13<<16 + 11<<12
	NotFound             = Storage + 1<<8 // Not found (subclass DB2XX)
	AlreadyExists        = Storage + 2<<8 // A resource with same index already exists (code DB101)
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
