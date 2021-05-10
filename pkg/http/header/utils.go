package header

import (
	"net/http"
)

func Copy(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func FromMap(m map[string][]string) http.Header {
	header := make(http.Header)
	for k, vv := range m {
		for _, v := range vv {
			header.Add(k, v)
		}
	}
	return header
}

func Overide(dst http.Header, overides map[string][]string) {
	for header, vv := range overides {
		if len(vv) == 0 {
			dst.Del(header)
		} else {
			for _, v := range vv {
				if v != "" {
					dst.Set(header, v)
				}
			}
		}
	}
}
