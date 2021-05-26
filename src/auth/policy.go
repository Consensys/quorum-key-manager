package auth

type Policy struct {
	// Name of the policy
	Name string

	// Statements of the Policy
	Statements []*Statement
}

type Statement struct {
	// Name of the statement
	Name string

	// Effect of the statement ('Allow' or 'Deny')
	Effect string

	// Actions the statement covers
	Actions []string

	// Resource the statement covers
	Resource []string
}
