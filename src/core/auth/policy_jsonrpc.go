package auth

import (
	"fmt"
	"strings"
)

const serviceMethodSeparator = "_"

type JSONRPCScope struct {
	Service string
	Method  string
}

// JSONRPCEndorsement JSONRPC accessed
type JSONRPCEndorsement struct {
	// Methods should be exact methods or some method patterns (e.g. "eth_sendTransaction", "eth_*", "*")
	Scopes []*JSONRPCScope
}

// IsAuthorized indicates wether
func (e *JSONRPCEndorsement) IsAuthorized(method string) error {
	var service string
	elems := strings.SplitN(method, serviceMethodSeparator, 2)
	if len(elems) == 2 {
		service = elems[0]
		method = elems[1]
	}

	for _, scope := range e.Scopes {
		if scope.Service == "*" && scope.Method == "*" {
			return nil
		}
		if scope.Service == "*" && scope.Method == method {
			return nil
		}
		if scope.Service == service && scope.Method == "*" {
			return nil
		}
		if scope.Service == service && scope.Method == method {
			return nil
		}
	}

	return fmt.Errorf("not authorized error")
}

func (policies *Policies) IsJSONRPCAuthorized(method string) error {
	for _, plcy := range (*policies)[PolicyTypeJSONRPC] {
		if err := plcy.Endorsement.(*JSONRPCEndorsement).IsAuthorized(method); err == nil {
			return nil
		}
	}

	return fmt.Errorf("not authorized")
}
