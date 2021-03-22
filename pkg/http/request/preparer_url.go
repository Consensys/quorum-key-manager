package request

import (
	"net/http"
	"net/url"
)

// URL set request URL to the passed URL
func URL(u *url.URL) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		newReq := *req
		newReq.URL = CopyURL(u)
		return &newReq, nil
	})
}

func CopyURL(i *url.URL) *url.URL {
	out := *i
	if i.User != nil {
		out.User = new(url.Userinfo)
		*out.User = *i.User
	}
	return &out
}
