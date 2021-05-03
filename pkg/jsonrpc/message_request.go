package jsonrpc

import (
	"context"
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

	ctx context.Context

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

func (msg *RequestMsg) Context() context.Context {
	if msg.ctx == nil {
		return context.Background()
	}
	return msg.ctx
}

func (msg *RequestMsg) WithContext(ctx context.Context) *RequestMsg {
	cpy := msg.Copy()
	cpy.ctx = ctx
	return cpy
}

func (msg *RequestMsg) Copy() *RequestMsg {
	newMsg := new(RequestMsg)

	newMsg.Version = msg.Version
	newMsg.Method = msg.Method
	newMsg.Params = msg.Params
	newMsg.ID = msg.ID

	if msg.raw != nil {
		newMsg.raw = new(jsonReqMsg)
		newMsg.raw.Version = msg.raw.Version
		newMsg.raw.Method = msg.raw.Method

		if msg.raw.ID != nil {
			newID := make([]byte, len(*msg.raw.ID))
			newMsg.raw.ID = (*json.RawMessage)(&newID)
			copy(newID, *msg.raw.ID)
		}

		if msg.raw.Params != nil {
			newParams := make([]byte, len(*msg.raw.Params))
			newMsg.raw.Params = (*json.RawMessage)(&newParams)
			copy(newParams, *msg.raw.Params)
		}
	}

	return newMsg
}

// Validate JSON-Requests body
func (msg *RequestMsg) Validate() error {
	if msg.Version == "" {
		return InvalidRequest(fmt.Errorf("missing version"))
	}

	if msg.Method == "" {
		return InvalidRequest(fmt.Errorf("missing method"))
	}

	if msg.ID != nil {
		err := validateID(msg.ID)
		if err != nil {
			return InvalidRequest(err)
		}
	}

	return nil
}

// UnmarshalID into v
func (msg *RequestMsg) UnmarshalID(v interface{}) error {
	var err error
	if msg.raw != nil && msg.raw.ID != nil {
		err = json.Unmarshal(*msg.raw.ID, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithID(reflect.ValueOf(v).Elem().Interface())
	}

	return err
}

// UnmarshalParams into v
func (msg *RequestMsg) UnmarshalParams(v interface{}) error {
	if msg.raw == nil {
		return fmt.Errorf("cannot unmarshal params from a non unmarshaled request message")
	}

	var err error
	if msg.raw.Params != nil {
		err = json.Unmarshal(*msg.raw.Params, v)
	} else {
		err = json.Unmarshal(null, v)
	}

	if err == nil {
		msg.WithParams(reflect.ValueOf(v).Elem().Interface())
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
