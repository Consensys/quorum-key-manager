package response

import (
	"compress/gzip"
	"net/http"
)

func GZIP() Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			resp.Body = reader
		}
		return nil
	})
}
