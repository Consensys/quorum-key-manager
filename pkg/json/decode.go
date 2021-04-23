package json

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
)

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return err
	}

	err = getValidator().Struct(req)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errMessage string
			for _, fe := range ves {
				errMessage += fmt.Sprintf("field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}

			return fmt.Errorf(errMessage)
		}

		return fmt.Errorf("invalid body")
	}

	return nil
}
