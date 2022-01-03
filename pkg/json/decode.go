package json

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/go-playground/validator/v10"
)

func UnmarshalBody(body io.Reader, req interface{}) error {
	if body == nil {
		return errors.InvalidFormatError("body is nil")
	}
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return err
	}

	return validateStruct(req)
}

func UnmarshalJSON(src, dest interface{}) error {
	bdata, err := MarshalJSON(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bdata, dest)
	if err != nil {
		return err
	}

	return validateStruct(dest)
}

func UnmarshalYAML(src, dest interface{}) error {
	bdata, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bdata, dest)
	if err != nil {
		return err
	}

	return validateStruct(dest)
}

func validateStruct(s interface{}) error {
	err := getValidator().Struct(s)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errMessage string
			for _, fe := range ves {
				errMessage += fmt.Sprintf("field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}

			return fmt.Errorf(errMessage)
		}

		return err
	}

	return nil
}
