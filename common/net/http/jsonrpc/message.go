package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var null = json.RawMessage("null")

// RequestMsg is a struct allowing to encode/decode a JSON-RPC request body
type RequestMsg struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Validate JSON-Requests body
func (msg *RequestMsg) Validate() error {
	if msg.Version == "" {
		return fmt.Errorf("missing version")
	}

	if msg.Method == "" {
		return fmt.Errorf("missing method")
	}

	err := validateID(msg.ID)
	if err != nil {
		return err
	}

	return nil
}

func validateID(id []byte) error {
	if len(id) == 0 {
		return fmt.Errorf("missing id")
	}

	if id[0] != '{' && id[0] != '[' {
		return nil
	}

	return fmt.Errorf("invalid id %v", string(id))
}

// WithID attaches ID
func (msg *RequestMsg) WithID(id interface{}) error {
	if id == nil {
		msg.ID = null
		return nil
	}

	b, err := json.Marshal(id)
	if err != nil {
		return err
	}

	err = validateID(b)
	if err != nil {
		return err
	}

	msg.ID = b

	return nil
}

// WithParams attaches parameters
func (msg *RequestMsg) WithParams(params interface{}) error {
	if params == nil {
		msg.Params = null
		return nil
	}

	var err error
	if msg.Params, err = json.Marshal(params); err != nil {
		return err
	}

	return nil
}

// UnmarshalParams into v
func (msg *RequestMsg) UnmarshalParams(v interface{}) error {
	return json.Unmarshal(msg.Params, v)
}

// ErrorMsg is a struct allowing to encode/decode a JSON-RPC response error
type ErrorMsg struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// WithData attaches data
func (msg *ErrorMsg) WithData(data interface{}) error {
	if data == nil {
		msg.Data = null
		return nil
	}

	var err error
	msg.Data, err = json.Marshal(data)

	return err
}

// UnmarshalData into v
func (msg *ErrorMsg) UnmarshalData(v interface{}) error {
	return json.Unmarshal(msg.Data, v)
}

// Error function to match the error interface
func (msg *ErrorMsg) Error() string {
	return msg.Message
}

// ResponseMsg is a struct allowing to encode/decode JSON-RPC response body
type ResponseMsg struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ErrorMsg       `json:"error,omitempty"`
}

// Validate JSON-RPC reseponse is valid
func (msg *ResponseMsg) Validate() error {
	if msg.Version == "" {
		return fmt.Errorf("missing version")
	}

	isSuccess := msg.Error == nil
	hasResult := !(len(msg.Result) == 0 || bytes.Equal(msg.Result, null))

	if isSuccess && !hasResult {
		return fmt.Errorf("missing result on success")
	}

	if !isSuccess && hasResult {
		return fmt.Errorf("non empty result on failure")
	}

	err := validateID(msg.ID)
	if err != nil {
		return err
	}

	return nil
}

// WithID attaches ID
func (msg *ResponseMsg) WithID(id interface{}) error {
	if id == nil {
		msg.ID = null
		return nil
	}

	b, err := json.Marshal(id)
	if err != nil {
		return err
	}

	err = validateID(b)
	if err != nil {
		return err
	}

	msg.ID = b

	return nil
}

// WithResult attaches result
func (msg *ResponseMsg) WithResult(result interface{}) error {
	if result == nil {
		msg.Result = null
	}

	var err error
	msg.Result, err = json.Marshal(result)

	return err
}

// WithError attaches error
func (msg *ResponseMsg) WithError(err error) {
	if err == nil {
		msg.Error = nil
		return
	}

	if errMsg, ok := err.(*ErrorMsg); ok {
		msg.Error = errMsg
	} else {
		msg.Error = &ErrorMsg{
			Message: err.Error(),
		}
	}
}

// UnmarshalResult into v
func (msg *ResponseMsg) UnmarshalResult(v interface{}) error {
	return json.Unmarshal(msg.Result, v)
}
