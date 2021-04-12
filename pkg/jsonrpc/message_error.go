package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ErrorMsg is a struct allowing to manipulate a JSON-RPC response error
type ErrorMsg struct {
	Code    int
	Message string
	Data    interface{}

	raw *jsonErrMsg
}

// jsonRespMsg is a struct allowing to encode/decode a JSON-RPC response body
type jsonErrMsg struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

func (msg *ErrorMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonErrMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	msg.raw = raw
	msg.Code = raw.Code
	msg.Message = raw.Message

	if raw.Data != nil {
		msg.Data = *raw.Data
	}

	return nil
}

// MarshalJSON
func (msg *ErrorMsg) MarshalJSON() ([]byte, error) {
	raw := new(jsonErrMsg)

	raw.Code = msg.Code
	raw.Message = msg.Message

	raw.Data = new(json.RawMessage)
	if msg.Data != nil {
		b, err := json.Marshal(msg.Data)
		if err != nil {
			return nil, err
		}
		*raw.Data = b
	} else {
		copy(*raw.Data, null)
	}

	return json.Marshal(raw)
}

// WithData attaches data
func (msg *ErrorMsg) WithData(data interface{}) *ErrorMsg {
	msg.Data = data
	return msg
}

// UnmarshalData into v
func (msg *ErrorMsg) UnmarshalData(v interface{}) error {
	var err error
	if msg.raw.Data != nil {
		err = json.Unmarshal(*msg.raw.Data, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		_ = msg.WithData(v)
	}

	return nil
}

// Error function to match the error interface
func (msg *ErrorMsg) Error() string {
	return msg.Message
}

var jsonMessageType = reflect.TypeOf(json.RawMessage(nil))

func validateID(id interface{}) error {
	idV := reflect.ValueOf(id)
	if idV.IsZero() {
		return nil
	}

	if idV.Type() == jsonMessageType {
		return validateRawID(idV.Interface().(json.RawMessage))
	}

	switch idV.Kind() {
	case reflect.Int, reflect.String:
		return nil
	case reflect.Ptr:
		return validateID(idV.Elem().Interface())
	default:
		return fmt.Errorf("invalid id (should be int or string but got %T)", id)
	}
}

func validateRawID(id json.RawMessage) error {
	if len(id) > 0 && id[0] != '{' && id[0] != '[' {
		return nil
	}

	return fmt.Errorf("invalid id %v", string(id))
}
