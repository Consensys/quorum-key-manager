package types

import (
	"net/http"
)

// UserContext is a set of data attached to every incoming request
type UserContext struct {
	// UserInfo records user information
	UserInfo *UserInfo
}

func NewUserContext(req *http.Request) *UserContext {
	return &UserContext{}
}
