package aliasent

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored
// in Value. The Value has 2 formats:
// - a JSON string if AliasKind is an AliasKindString.
// - a JSON array of strings if AliasKind is an AliasKindArray.
type Alias struct {
	Key          AliasKey
	RegistryName RegistryName
	// Value is an array containing Tessera/Orion keys base64 encoded in strings.
	Value AliasValue
}

type AliasKey string

type AliasValue []string

type RegistryName string
