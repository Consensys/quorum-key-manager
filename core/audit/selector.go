package audit

type Operator int

const (
	DoesNotExist Operator = iota
	Equals
	DoubleEquals
	In
	NotEquals
	NotIn
	Exists
)

var Operators = map[string]Operator{
	"!":      DoesNotExist,
	"=":      Equals,
	"==":     DoubleEquals,
	"in":     In,
	"!=":     NotEquals,
	"notin":  NotIn,
	"exists": Exists,
}

// Selector allows filters audit operations that matches some conditions
type Selector struct {
	// Type of the operation
	Type string `json:"type,omitempty"`

	// Data on the operation
	Data []*DataRequirement `json:"data,omitempty"`
}

type DataRequirement struct {
	Key      string   `json:"key,omitempty"`
	Operator Operator `json:"operator,omitempty"`
	Values   []string `json:"values,omitempty"`
}
