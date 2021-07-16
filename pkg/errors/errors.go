package errors

import (
	"fmt"
)

type Error struct {
	Message string
	Code    string
}

func (e *Error) GetCode() string {
	return e.Code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) GetMessage() string {
	return e.Message
}

func Errorf(code, format string, a ...interface{}) *Error {
	return &Error{fmt.Sprintf(format, a...), code}
}

func FromError(err interface{}) *Error {
	if err == nil {
		return nil
	}

	ierr, ok := err.(*Error)
	if !ok {
		return Errorf(Internal, err.(error).Error())
	}

	return ierr
}
