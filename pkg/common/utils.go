package common

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
