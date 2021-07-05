package response

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/http/header"
)

func HeadersModifier(h func(http.Header) error) Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		return h(resp.Header)
	})
}

// Headers sets or deletes custom request headers
func Headers(overides map[string][]string) Modifier {
	return HeadersModifier(header.Overide(overides))
}
