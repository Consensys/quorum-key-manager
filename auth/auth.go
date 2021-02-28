package auth

// Auth is the piece information attached to every client request after
// authentication succeeded (e.g. JWT based auth, Mutual TLS auth, etc.)

// It shoud be computed by an authentication middleware responsible to validate authentication
// and then attach Auth to the request context so it can flow to sub-sequent components
// that can perform authorization/permissions checks (e.g. JSON-RPC methods checks, store access check, etc.)
type Auth struct {
	// ID of the authenticated client
	ID string

	// Policies associated to the authenticated client
	Policies []string

	// Metadata is arbitrary string-type metadata
	Metadata map[string]string

	// Raw is JSON-encodable data that is stored with the auth struct (e.g. JWT Token)
	Raw map[string]interface{}
}
