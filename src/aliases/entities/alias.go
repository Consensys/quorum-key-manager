package entities

// Alias allows the user to associates a RegistryName + a Key to 1 or more
// public keys stored in Value.
type Alias struct {
	// Key is the unique alias key.
	Key string

	// RegistryName is the unique registry name.
	RegistryName string

	// Value is a slice containing Tessera/Orion keys base64 encoded in strings.
	Value AliasValue
}

type AliasValueInterface interface {
	Kind() Kind
	Value() interface{}
}

type AliasValue struct {
	Kind  Kind
	Value interface{}
}

type Kind string

const (
	KindUnknown Kind = ""
	KindString       = "string"
	KindArray        = "array"
)

type StringValue string

func (sv StringValue) Kind() Kind {
	return KindString
}

func (sv StringValue) Value() interface{} {
	return sv
}

type ArrayValue string

func (av ArrayValue) Kind() Kind {
	return KindString
}

func (av ArrayValue) Value() interface{} {
	return av
}
