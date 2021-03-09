package json

import (
	"encoding/json"
	"fmt"
	"io"
)

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return fmt.Errorf("failed to decode request body")
	}

	return nil
}
