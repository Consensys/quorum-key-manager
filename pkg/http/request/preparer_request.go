package request

import (
	"net/http"
	"net/url"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/header"
)

func OverideURL(dst, src *url.URL) {
	if src.Scheme != "" {
		dst.Scheme = src.Scheme
	}

	if src.Opaque != "" {
		dst.Opaque = src.Opaque
	}

	if src.User != nil {
		dst.User = src.User
	}

	if src.Host != "" {
		dst.Host = src.Host
	}

	if src.Path != "" {
		dst.Path = src.Path
	}

	if src.RawPath != "" {
		dst.RawPath = src.RawPath
	}

	if src.ForceQuery {
		dst.ForceQuery = src.ForceQuery
	}

	if src.RawQuery != "" {
		dst.RawQuery = src.RawQuery
	}

	if src.Fragment != "" {
		dst.Fragment = src.Fragment
	}

	if src.RawFragment != "" {
		dst.RawFragment = src.RawFragment
	}
}

// Request is a preparer that enhance req with baseReq fields
func Request(baseReq *http.Request) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		// Clone from base request
		outReq := baseReq.Clone(req.Context())

		// Set Method
		outReq.Method = req.Method

		// Overide Header
		err := header.Overide(req.Header)(outReq.Header)
		if err != nil {
			return nil, err
		}

		// Overide URL
		if outReq.URL == nil {
			outReq.URL = CopyURL(req.URL)
		} else {
			OverideURL(outReq.URL, req.URL)
		}

		// Overide URI
		if req.RequestURI != "" {
			outReq.RequestURI = req.RequestURI
		}

		// Overide body
		outReq.Body = req.Body
		outReq.GetBody = req.GetBody
		outReq.ContentLength = req.ContentLength

		return outReq, nil
	})
}
