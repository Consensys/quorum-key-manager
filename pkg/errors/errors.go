//nolint
package errors

import (
	"fmt"
)

type Error struct {
	Message   string
	Code      uint64
	Component string
}

func (e *Error) GetCode() uint64 {
	return e.Code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s", e.Component, e.Message)
}

func (m *Error) GetMessage() string {
	return m.Message
}

func Errorf(code uint64, format string, a ...interface{}) *Error {
	return &Error{fmt.Sprintf(format, a...), code, ""}
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
