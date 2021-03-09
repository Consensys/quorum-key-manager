package entities

// StoreInfo for a store
type StoreInfo struct {
	// Kind of store
	Kind string

	// Name set by user
	Name string

	// Info about the store proper to each implementation
	// It should not expose any secret information about the store configuration
	Info interface{}

	// Tags set by user when creating the stores
	Tags map[string]string
}
