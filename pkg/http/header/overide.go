package header

import "net/http"

func Override(overrides map[string][]string) func(http.Header) error {
	return func(dst http.Header) error {
		for header, vv := range overrides {
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

		return nil
	}
}
