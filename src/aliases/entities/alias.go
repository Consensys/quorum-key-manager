package entities

// Alias allows the user to associates a RegistryName + a Key to 1 or more
// public keys stored in Value.
type Alias struct {
	// Key is the unique alias key.
	Key string

	// RegistryName is the unique registry name.
	RegistryName string

	// Value is a slice containing Tessera/Orion keys base64 encoded in strings.
	Value []string
}
