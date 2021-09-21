package jsonrpc

import (
	"context"
	"fmt"
	"reflect"
	"unicode"
)

var (
	errorType   = reflect.TypeOf(new(error)).Elem()
	contextType = reflect.TypeOf(new(context.Context)).Elem()
	clientType  = reflect.TypeOf(new(Client)).Elem()
)

// ProvideCaller takes a list of user defines callers as argument and
// automatically populates all caller's fields with RPC functions
//
// It aims at facilitate the implementation of Web3 client connecting to downstream node
//
// - Caller MUST be pointers to struct
// - All caller's fields MUST be functions mathc
// - Caller field func MUST accept a single input of type jsonrpc.Client and MUST return a single output which is a function
// - Caller field output func MUST return at most 2 outputs (if 2 the second MUST be an error)
//
// Example of valid caller struct:
//
// type ExampleCaller struct {
// 	   CtxInput_NoOutput        func(Client) func(context.Context)
// 	   NoInput_NoOutput         func(Client) func()
// 	   NonCtxInput_NoOutput     func(Client) func(int)
// 	   MultiInput_NoOutput      func(Client) func(context.Context, int, string)
// 	   NoInput_ErrorOutput      func(Client) func() error
// 	   NoInput_IntOutput        func(Client) func() int
// 	   NoInput_IntErrorOutput   func(Client) func() (int, error)
// 	   StructInput_StructOutput func(Client) func(context.Context, *TestParam) (*TestResult, error)
// 	   AllTags                  func(Client) func()                                                 `method:"exampleMethod" namespace:"eth"`
// 	   MethodTag                func(Client) func()                                                 `method:"exampleMethod"`
// 	   NamespaceTag             func(Client) func()                                                 `namespace:"eth"`
// 	   ObjectTag                func(Client) func(*TestParam)                                       `object:"-"`
// }
func ProvideCaller(callers ...interface{}) error {
	for _, caller := range callers {
		cllrTyp := reflect.TypeOf(caller)
		if cllrTyp.Kind() != reflect.Ptr || cllrTyp.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("caller must be a pointer to a struct")
		}

		cllrVal := reflect.ValueOf(caller)

		for i := 0; i < cllrTyp.Elem().NumField(); i++ {
			field := cllrTyp.Elem().Field(i)
			fn, err := makeRPCCallerFunc(&field)
			if err != nil {
				return err
			}

			cllrVal.Elem().Field(i).Set(fn)
		}
	}

	return nil
}

func makeRPCCallerFunc(f *reflect.StructField) (reflect.Value, error) {
	ftyp := f.Type
	if ftyp.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("caller's fields must be functions")
	}

	if ftyp.NumIn() != 1 && ftyp.In(0) != clientType {
		return reflect.Value{}, fmt.Errorf("caller's field func must accept a single input of type %v", clientType)
	}

	if ftyp.NumOut() != 1 || ftyp.Out(0).Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("caller's field func must return a single output which is a function")
	}

	fun := &rpcCallerFunc{
		ftyp:  ftyp.Out(0),
		retry: f.Tag.Get("retry") == "true",
	}

	method, hasMethod := f.Tag.Lookup("method")
	namespace, hasNamespace := f.Tag.Lookup("namespace")

	switch {
	case hasMethod && hasNamespace:
		fun.method = fmt.Sprintf("%v_%v", namespace, method)
	case hasMethod:
		fun.method = method
	case hasNamespace:
		fun.method = fmt.Sprintf("%v_%v", namespace, formatName(f.Name))
	default:
		fun.method = f.Name
	}

	_, fun.object = f.Tag.Lookup("object")

	if fun.ftyp.NumIn() > 0 && fun.ftyp.In(0) == contextType {
		fun.hasCtx = 1
	}

	var err error
	fun.valOut, fun.errOut, fun.nout, err = processFuncOut(fun.ftyp)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.MakeFunc(ftyp, func(callArg []reflect.Value) []reflect.Value {
		return []reflect.Value{
			reflect.MakeFunc(ftyp.Out(0), func(args []reflect.Value) []reflect.Value {
				return fun.handleCall(callArg[0].Interface().(Client), args...)
			}),
		}
	}), nil
}

// processFuncOut finds value and error Outs in function
func processFuncOut(funcType reflect.Type) (valOut, errOut, n int, err error) {
	errOut = -1 // -1 if not found
	valOut = -1
	n = funcType.NumOut()

	switch n {
	case 0:
	case 1:
		if funcType.Out(0) == errorType {
			errOut = 0
		} else {
			valOut = 0
		}
	case 2:
		valOut = 0
		errOut = 1
		if funcType.Out(1) != errorType {
			err = fmt.Errorf("caller's field output function second output must be an error")
		}
	default:
		err = fmt.Errorf("caller's field output function must return at most 2 outputs")
	}

	return
}

type rpcCallerFunc struct {
	ftyp   reflect.Type
	method string

	nout   int
	valOut int
	errOut int
	hasCtx int

	object bool
	retry  bool
}

func (fn *rpcCallerFunc) prepareParams(args ...reflect.Value) (params interface{}) {
	switch {
	case len(args) == 1 && fn.object:
		params = args[0].Interface()
	default:
		prms := make([]interface{}, len(args))
		for i, arg := range args {
			prms[i] = arg.Interface()
		}
		params = prms
	}
	return
}

func (fn *rpcCallerFunc) handleCall(client Client, args ...reflect.Value) []reflect.Value {
	params := fn.prepareParams(args[fn.hasCtx:]...)

	var ctx context.Context
	if fn.hasCtx == 1 {
		ctx = args[0].Interface().(context.Context)
	} else {
		ctx = context.Background()
	}

	resp, err := client.Do((&RequestMsg{}).WithContext(ctx).WithMethod(fn.method).WithParams(params))
	if err != nil {
		return fn.processError(err)
	}

	return fn.processResponse(resp)
}

func (fn *rpcCallerFunc) processError(err error) []reflect.Value {
	out := make([]reflect.Value, fn.nout)

	if fn.valOut != -1 {
		out[fn.valOut] = reflect.New(fn.ftyp.Out(fn.valOut)).Elem()
	}

	if fn.errOut != -1 {
		out[fn.errOut] = reflect.ValueOf(err)
	}

	return out
}

func (fn *rpcCallerFunc) processResponse(resp *ResponseMsg) []reflect.Value {
	out := make([]reflect.Value, fn.nout)

	if fn.valOut != -1 {
		val := reflect.New(fn.ftyp.Out(fn.valOut))
		if err := resp.UnmarshalResult(val.Interface()); err != nil && resp.Err() == nil {
			// We exit only if there is not json-rpc error to parse
			return fn.processError(err)
		}
		out[fn.valOut] = val.Elem()
	}

	if fn.errOut != -1 {
		if resp.Err() != nil {
			out[fn.errOut] = reflect.ValueOf(resp.Err())
		} else {
			out[fn.errOut] = reflect.New(errorType).Elem()
		}

	}

	return out
}

// formatName converts to first character of name to lowercase.
func formatName(name string) string {
	ret := []rune(name)
	if len(ret) > 0 {
		ret[0] = unicode.ToLower(ret[0])
	}
	return string(ret)
}
