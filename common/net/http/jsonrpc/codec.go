package jsonrpc

import (
	"encoding/json"
	"io"
)

type ClientCodec interface {
	WriteRequest(msg *RequestMsg) error
	ReadResponse(msg *ResponseMsg) error
}

type clientCodec struct {
	respdec *json.Decoder
	reqenc  *json.Encoder
}

func newClientCodec(reqw io.Writer, respr io.Reader) *clientCodec {
	return &clientCodec{
		respdec: json.NewDecoder(respr),
		reqenc:  json.NewEncoder(reqw),
	}
}

func (codec *clientCodec) WriteRequest(msg *RequestMsg) error {
	return codec.reqenc.Encode(msg)
}

func (codec *clientCodec) ReadResponse(msg *ResponseMsg) error {
	return codec.respdec.Decode(msg)
}

type ServerCodec interface {
	ReadRequest(msg *RequestMsg) error
	WriteResponse(msg *ResponseMsg) error
}

type serverCodec struct {
	reqdec  *json.Decoder
	respenc *json.Encoder
}

func newServerCodec(reqr io.Reader, respw io.Writer) *serverCodec {
	return &serverCodec{
		reqdec:  json.NewDecoder(reqr),
		respenc: json.NewEncoder(respw),
	}
}

func (codec *serverCodec) ReadRequest(msg *RequestMsg) error {
	return codec.reqdec.Decode(msg)
}

func (codec *serverCodec) WriteResponse(msg *ResponseMsg) error {
	return codec.respenc.Encode(msg)
}
