package request

import "net/http"

func Request(baseReq *http.Request) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		// Clone from base request
		outReq := baseReq.Clone(req.Context())

		// Set Method
		outReq.Method = req.Method

		for k := range req.Header {
			outReq.Header.Set(k, req.Header.Get(k))
		}

		if outReq.URL == nil {
			outReq.URL = CopyURL(req.URL)
		}

		if req.URL.Scheme != "" {
			outReq.URL.Scheme = req.URL.Scheme
		}

		if req.URL.Opaque != "" {
			outReq.URL.Opaque = req.URL.Opaque
		}

		if req.URL.User != nil {
			outReq.URL.User = req.URL.User
		}

		if req.URL.Host != "" {
			outReq.URL.Host = req.URL.Host
		}

		if req.URL.Path != "" {
			outReq.URL.Path = req.URL.Path
		}

		if req.URL.RawPath != "" {
			outReq.URL.RawPath = req.URL.RawPath
		}

		if req.URL.ForceQuery {
			outReq.URL.ForceQuery = req.URL.ForceQuery
		}

		if req.URL.RawQuery != "" {
			outReq.URL.RawQuery = req.URL.RawQuery
		}

		if req.URL.Fragment != "" {
			outReq.URL.Fragment = req.URL.Fragment
		}

		if req.URL.RawFragment != "" {
			outReq.URL.RawFragment = req.URL.RawFragment
		}

		if req.RequestURI != "" {
			outReq.RequestURI = req.RequestURI
		}

		outReq.Body = req.Body
		outReq.GetBody = req.GetBody
		outReq.ContentLength = req.ContentLength

		return outReq, nil
	})
}
