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
	callerType  = reflect.TypeOf(new(Caller)).Elem()
)

// Provide takes services as argument and populates fields with RPC functions
// It aims at facilitate the implemention of Web3 client connecting to downstream node

// - Services MUST be pointers to struct
// - All service's fields MUST be functions mathc
// - Service field func MUST accept a single input of type jsonrpc.Caller and MUST return a single output which is a function
// - Service field output func MUST return at most 2 outputs (if 2 the second MUST be an error)

// Example of valid service struct:

// type ExampleService struct {
// 	   CtxInput_NoOutput        func(Caller) func(context.Context)
// 	   NoInput_NoOutput         func(Caller) func()
// 	   NonCtxInput_NoOutput     func(Caller) func(int)
// 	   MultiInput_NoOutput      func(Caller) func(context.Context, int, string)
// 	   NoInput_ErrorOutput      func(Caller) func() error
// 	   NoInput_IntOutput        func(Caller) func() int
// 	   NoInput_IntErrorOutput   func(Caller) func() (int, error)
// 	   StructInput_StructOutput func(Caller) func(context.Context, *TestParam) (*TestResult, error)
// 	   AllTags                  func(Caller) func()                                                 `method:"exampleMethod" namespace:"eth"`
// 	   MethodTag                func(Caller) func()                                                 `method:"exampleMethod"`
// 	   NamespaceTag             func(Caller) func()                                                 `namespace:"eth"`
// 	   ObjectTag                func(Caller) func(*TestParam)                                       `object:"-"`
// }
func Provide(services ...interface{}) error {
	for _, service := range services {
		srvTyp := reflect.TypeOf(service)
		if srvTyp.Kind() != reflect.Ptr || srvTyp.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("service must be a pointer to a struct")
		}

		srvVal := reflect.ValueOf(service)

		for i := 0; i < srvTyp.Elem().NumField(); i++ {
			field := srvTyp.Elem().Field(i)
			fn, err := makeRPCFunc(&field)
			if err != nil {
				return err
			}

			srvVal.Elem().Field(i).Set(fn)
		}
	}

	return nil
}

func makeRPCFunc(f *reflect.StructField) (reflect.Value, error) {
	ftyp := f.Type
	if ftyp.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("service's field must be a function")
	}

	if ftyp.NumIn() != 1 && ftyp.In(0) != contextType {
		return reflect.Value{}, fmt.Errorf("service field func must accept a single input of type %v", callerType)
	}

	if ftyp.NumOut() != 1 || ftyp.Out(0).Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("service field func must return a single output which is a function")
	}

	fun := &rpcFunc{
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
				return fun.handleCall(callArg[0].Interface().(Caller), args...)
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
			err = fmt.Errorf("service field output function second output must be an error")
		}
	default:
		err = fmt.Errorf("service field output function must return at most 2 outputs")
	}

	return
}

type rpcFunc struct {
	ftyp   reflect.Type
	method string

	nout   int
	valOut int
	errOut int
	hasCtx int

	object bool
	retry  bool
}

func (fn *rpcFunc) prepareParams(args ...reflect.Value) (params interface{}) {
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

func (fn *rpcFunc) handleCall(cllr Caller, args ...reflect.Value) []reflect.Value {
	params := fn.prepareParams(args[fn.hasCtx:]...)

	var ctx context.Context
	if fn.hasCtx == 1 {
		ctx = args[0].Interface().(context.Context)
	} else {
		ctx = context.Background()
	}

	resp, err := cllr.Call(ctx, fn.method, params)
	if err != nil {
		return fn.processError(err)
	}

	return fn.processResponse(resp)
}

func (fn *rpcFunc) processError(err error) []reflect.Value {
	out := make([]reflect.Value, fn.nout)

	if fn.valOut != -1 {
		out[fn.valOut] = reflect.New(fn.ftyp.Out(fn.valOut)).Elem()
	}

	if fn.errOut != -1 {
		out[fn.errOut] = reflect.ValueOf(err)
	}

	return out
}

func (fn *rpcFunc) processResponse(resp *Response) []reflect.Value {
	out := make([]reflect.Value, fn.nout)

	if fn.valOut != -1 {
		val := reflect.New(fn.ftyp.Out(fn.valOut))
		if err := resp.UnmarshalResult(val.Interface()); err != nil {
			return fn.processError(err)
		}

		out[fn.valOut] = val.Elem()
	}

	if fn.errOut != -1 {
		out[fn.errOut] = reflect.New(errorType).Elem()
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
