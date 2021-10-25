package entities

// UserInfo are extracted from request credentials by authentication middleware
type UserInfo struct {
	// AuthMode records the mode that succeeded to Authenticate the request ('tls', 'api-key', 'oidc' or '')
	AuthMode string

	// Tenant belonged by the user
	Tenant string

	// Subject identifies the user
	Username string

	// Roles indicates the user's membership
	Roles []string

	// Permissions specify
	Permissions []Permission
}

var WildcardUser = &UserInfo{
	Permissions: ListPermissions(),
}

var AnonymousUser = &UserInfo{
	Username:    "anonymous",
	Roles:       []string{AnonymousRole},
	Permissions: []Permission{},
}
