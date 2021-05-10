package common

import "reflect"

func Tomapstrptr(m map[string]string) map[string]*string {
	nm := make(map[string]*string)
	for k, v := range m {
		nm[k] = &(&struct{ x string }{v}).x
	}
	return nm
}

func Tomapstr(m map[string]*string) map[string]string {
	nm := make(map[string]string)
	for k, v := range m {
		nm[k] = *v
	}
	return nm
}

func ToPtr(v interface{}) interface{} {
	p := reflect.New(reflect.TypeOf(v))
	p.Elem().Set(reflect.ValueOf(v))
	return p.Interface()
}
