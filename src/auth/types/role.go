package types

type Role struct {
	Name        string
	Permissions []Permission
}

const AnonymousRole = "anonymous"
