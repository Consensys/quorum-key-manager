package response

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/header"
)

// Headers sets or deletes custom request headers
func Headers(overides map[string][]string) Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		header.Overide(resp.Header, overides)
		return nil
	})
}
