package types

type Policy struct {
	// Name of the policy
	Name string

	// Statements of the Policy
	Statements []*Statement
}

type Statement struct {
	// Name of the statement
	Name string `json:"name"`

	// Effect of the statement ('Allow' or 'Deny')
	Effect string `json:"effect"`

	// Actions the statement covers
	Actions []string `json:"actions"`

	// Resource the statement covers
	Resource []string `json:"resource"`
}
