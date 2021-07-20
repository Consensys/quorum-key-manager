package types

// UserInfo are extracted from request credentials by authentication middleware
type UserInfo struct {
	// AuthMode records the mode that succeeded to Authenticate the request ('tls', 'api-key', 'oidc' or '')
	AuthMode string

	// Username identifies the user
	Username string

	// Groups indicates the user's membership to collection of users with specific permissions to access
	Groups []string
}

var AnonymousUser = &UserInfo{
	Username: "user:anonymous",
	Groups: []string{
		"system:unauthenticated",
	},
}

var AuthenticatedUser = &UserInfo{
	Username: "user:authenticated",
	Groups: []string{
		"system:authenticated",
	},
}
