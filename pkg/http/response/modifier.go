package response

import "net/http"

// Modifier is the interface that allows to modify http.http.Response the Respond method.
type Modifier interface {
	Modify(*http.Response) error
}

// ModifierFunc is a method that implements the Modifier interface.
type ModifierFunc func(*http.Response) error

// Respond implements the Modifier interface on ModifierFunc.
func (rf ModifierFunc) Modify(r *http.Response) error {
	return rf(r)
}

// CombineModifier combines multiple modifiers into a single one
func CombineModifier(modifiers ...Modifier) Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		var err error
		for _, modifier := range modifiers {
			err = modifier.Modify(resp)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

var NoopModifier = ModifierFunc(func(*http.Response) error { return nil })
