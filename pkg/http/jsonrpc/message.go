package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var null = json.RawMessage("null")

// RequestMsg allows to manipulate a JSON-RPC v2 request
type RequestMsg struct {
	Version string
	Method  string
	ID      interface{}
	Params  interface{}

	raw *jsonReqMsg
}

// jsonReqMsg is a struct allowing to encode/decode a JSON-RPC request body
type jsonReqMsg struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

// UnmarshalJSON
func (msg *RequestMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonReqMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	msg.raw = raw
	msg.Version = raw.Version
	msg.Method = raw.Method

	if raw.ID != nil {
		msg.ID = *raw.ID
	}

	if raw.Params != nil {
		msg.Params = *raw.Params
	}

	return nil
}

// MarshalJSON
func (msg *RequestMsg) MarshalJSON() ([]byte, error) {
	raw := new(jsonReqMsg)

	raw.Version = msg.Version
	raw.Method = msg.Method

	raw.ID = new(json.RawMessage)
	if msg.ID != nil {
		b, err := json.Marshal(msg.ID)
		if err != nil {
			return nil, err
		}

		*raw.ID = b
	} else {
		copy(*raw.ID, null)
	}

	raw.Params = new(json.RawMessage)
	if msg.Params != nil {
		b, err := json.Marshal(msg.Params)
		if err != nil {
			return nil, err
		}
		*raw.Params = b
	} else {
		copy(*raw.Params, null)
	}

	return json.Marshal(raw)
}

// Validate JSON-Requests body
func (msg *RequestMsg) Validate() error {
	if msg.Version == "" {
		return fmt.Errorf("missing version")
	}

	if msg.Method == "" {
		return fmt.Errorf("missing method")
	}

	if msg.ID != nil {
		err := validateID(msg.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalID into v
func (msg *RequestMsg) UnmarshalID(v interface{}) error {
	var err error
	if msg.raw.ID != nil {
		err = json.Unmarshal(*msg.raw.ID, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithID(v)
	}

	return err
}

// UnmarshalParams into v
func (msg *RequestMsg) UnmarshalParams(v interface{}) error {
	var err error
	if msg.raw.Params != nil {
		err = json.Unmarshal(*msg.raw.Params, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithParams(v)
	}

	return err
}

// WithVersion attaches version
func (msg *RequestMsg) WithVersion(v string) *RequestMsg {
	msg.Version = v
	return msg
}

// WithMethod attaches method
func (msg *RequestMsg) WithMethod(method string) *RequestMsg {
	msg.Method = method
	return msg
}

// WithID attaches ID
func (msg *RequestMsg) WithID(id interface{}) *RequestMsg {
	msg.ID = id
	return msg
}

// WithParams attaches parameters
func (msg *RequestMsg) WithParams(params interface{}) *RequestMsg {
	msg.Params = params
	return msg
}

// ResponseMsg allows to manipulate a JSON-RPC v2 response
type ResponseMsg struct {
	Version string
	ID      interface{}
	Result  interface{}
	Error   *ErrorMsg

	raw *jsonRespMsg
}

// jsonRespMsg is a struct allowing to encode/decode a JSON-RPC response body
type jsonRespMsg struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

func (msg *ResponseMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonRespMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	msg.raw = raw
	msg.Version = raw.Version

	if raw.ID != nil {
		msg.ID = *raw.ID
	}

	if raw.Result != nil {
		msg.Result = *raw.Result
	}

	if raw.Error != nil {
		msg.Error = new(ErrorMsg)
		err = json.Unmarshal(*raw.Error, msg.Error)
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON
func (msg *ResponseMsg) MarshalJSON() ([]byte, error) {
	raw := new(jsonRespMsg)

	raw.Version = msg.Version

	raw.ID = new(json.RawMessage)
	if msg.ID != nil {
		b, err := json.Marshal(msg.ID)
		if err != nil {
			return nil, err
		}
		*raw.ID = b
	} else {
		copy(*raw.ID, null)
	}

	raw.Result = new(json.RawMessage)
	if msg.Result != nil {
		b, err := json.Marshal(msg.Result)
		if err != nil {
			return nil, err
		}
		*raw.Result = b
	} else {
		copy(*raw.Result, null)
	}

	raw.Error = new(json.RawMessage)
	if msg.Error != nil {
		b, err := json.Marshal(msg.Error)
		if err != nil {
			return nil, err
		}
		*raw.Error = b
	} else {
		copy(*raw.Error, null)
	}

	return json.Marshal(raw)
}

// Validate JSON-RPC reseponse is valid
func (msg *ResponseMsg) Validate() error {
	if msg.Version == "" {
		return fmt.Errorf("missing version")
	}

	isSuccess := msg.Error == nil
	hasResult := msg.Result != nil

	if isSuccess && !hasResult {
		return fmt.Errorf("missing result on success")
	}

	if !isSuccess && hasResult {
		return fmt.Errorf("non empty result on failure")
	}

	if msg.ID != nil {
		err := validateID(msg.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// WithVersion attaches version
func (msg *ResponseMsg) WithVersion(v string) *ResponseMsg {
	msg.Version = v
	return msg
}

// WithID attaches ID
func (msg *ResponseMsg) WithID(id interface{}) *ResponseMsg {
	msg.ID = id
	return msg
}

// WithResult attaches result
func (msg *ResponseMsg) WithResult(result interface{}) *ResponseMsg {
	msg.Result = result
	return msg
}

// WithError attaches error
func (msg *ResponseMsg) WithError(err error) *ResponseMsg {
	if err == nil {
		msg.Error = nil
		return msg
	}

	if errMsg, ok := err.(*ErrorMsg); ok {
		msg.Error = errMsg
	} else {
		msg.Error = &ErrorMsg{
			Message: err.Error(),
		}
	}

	return msg
}

// UnmarshalResult into v
func (msg *ResponseMsg) UnmarshalResult(v interface{}) error {
	var err error
	if msg.raw.Result != nil {
		err = json.Unmarshal(*msg.raw.Result, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithResult(v)
	}

	return err
}

// UnmarshalID into v
func (msg *ResponseMsg) UnmarshalID(v interface{}) error {
	var err error
	if msg.raw.ID != nil {
		err = json.Unmarshal(*msg.raw.ID, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithID(v)
	}

	return err
}

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
