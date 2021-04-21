package jsonrpc

import (
	"encoding/json"
	"fmt"
)

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

func (msg *ResponseMsg) Err() error {
	if msg.Error == nil {
		return nil
	}

	return msg.Error
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
	err := validateID(id)
	if err != nil {
		panic(err)
	}
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
		msg.Error = Error(err)
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
	if msg.raw != nil && msg.raw.ID != nil {
		err = json.Unmarshal(*msg.raw.ID, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithID(v)
	}

	return err
}
