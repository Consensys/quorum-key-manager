package auth

// Auth is the information attached to every client request after authentication succeeded (e.g. JWT based auth, Mutual TLS auth, etc.)

// This is attached to context so it can flow to every core components to base to perform 
// authentication checks (e.g. JSON-RPC methods checks, store access check, etc.)
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
