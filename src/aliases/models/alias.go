package aliasmodels

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored
// in Value. The Value has 2 formats:
// - a JSON string if AliasKind is an AliasKindString.
// - a JSON array of strings if AliasKind is an AliasKindArray.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	Key          AliasKey     `pg:",pk"`
	RegistryName RegistryName `pg:",pk"`
	Kind         AliasKind
	Value        AliasValue
}

type AliasKey string

type AliasValue string

type AliasKind string

const (
	AliasKindUnknown = ""
	AliasKindString  = "string"
	AliasKindArray   = "array"
)

type RegistryName string
