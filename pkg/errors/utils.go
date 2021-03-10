package errors

import (
	"fmt"
)

func CombineErrors(errs ...error) error {
	var err error
	for _, e := range errs {
		if e == nil {
			continue
		}

		if err == nil {
			err = e
		} else {
			err = fmt.Errorf("%v; %v", e, err)
		}
	}
	return err
}
