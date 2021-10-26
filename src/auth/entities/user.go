package entities

// UserClaims represent raw claims extracted from an authentication method
type UserClaims struct {
	Subject string
	Scope   string
	Roles   string
}

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

func NewWildcardUser() *UserInfo {
	return &UserInfo{
		Permissions: ListPermissions(),
	}
}

func NewAnonymousUser() *UserInfo {
	return &UserInfo{
		Username:    "anonymous",
		Roles:       []string{AnonymousRole},
		Permissions: []Permission{},
	}
}
