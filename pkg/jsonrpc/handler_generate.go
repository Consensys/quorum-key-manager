package jsonrpc

import (
	"fmt"
	"reflect"
)

// MakeHandler takes a function and transform it into a jsonrpc.Handler
func MakeHandler(f interface{}) (Handler, error) {
	fVal := reflect.ValueOf(f)
	if !fVal.IsValid() {
		return nil, fmt.Errorf("can not generate handler from zero value")
	}

	ftyp := fVal.Type()
	if ftyp.Kind() != reflect.Func {
		return nil, fmt.Errorf("expect function but got %T", f)
	}

	h := &rpcHandler{
		f: fVal,
	}

	if numIn := ftyp.NumIn(); numIn > 0 {
		if ftyp.In(0) == contextType {
			h.hasCtx = 1
		}

		for i := h.hasCtx; i < numIn; i++ {
			h.paramsType = append(h.paramsType, ftyp.In(i))
		}
	}

	numOut := ftyp.NumOut()
	switch numOut {
	case 0:
		return nil, fmt.Errorf("function must return at least one output")
	case 1:
	case 2:
		if ftyp.Out(1) != errorType {
			return nil, fmt.Errorf("function second output must be an error")
		}
		h.hasError = true
	default:
		return nil, fmt.Errorf("function must return at most two outputs")
	}

	return h, nil
}

type rpcHandler struct {
	f reflect.Value

	paramsType []reflect.Type
	hasCtx     int
	hasError   bool
}

func (fn *rpcHandler) ServeRPC(rw ResponseWriter, msg *RequestMsg) {
	// Prepare params
	inParams := fn.newParams()

	params := prepareParams(inParams...)
	err := msg.UnmarshalParams(&params)
	if err != nil {
		_ = WriteError(rw, InvalidParamsError(err))
		return
	}

	var in []reflect.Value
	if fn.hasCtx > 0 {
		in = append(in, reflect.ValueOf(msg.Context()))
	}

	for _, inParam := range inParams {
		in = append(in, inParam.Elem())
	}

	out := fn.f.Call(in)

	if fn.hasError && !out[1].IsNil() {
		_ = WriteError(rw, out[1].Interface().(error))
		return
	}

	_ = WriteResult(rw, out[0].Interface())
}

func (fn *rpcHandler) newParams() []reflect.Value {
	var inputs []reflect.Value

	for _, typ := range fn.paramsType {
		inputs = append(inputs, reflect.New(typ))
	}

	return inputs
}

func prepareParams(inputs ...reflect.Value) []interface{} {
	params := []interface{}{}
	for _, in := range inputs {
		params = append(params, in.Interface())
	}

	return params
}
