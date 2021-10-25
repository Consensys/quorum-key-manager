package jwt

import "context"

type CustomClaims struct {
	Roles string `json:"roles"`
	Scope string `json:"scope"`
}

func (claims *CustomClaims) Validate(_ context.Context) error {
	// TODO: Apply validation on custom claims if needed, currently no validation is needed
	return nil
}
