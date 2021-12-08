package entities

const (
	EthereumStoreType = "ethereum"
	KeyStoreType      = "key"
	SecretStoreType   = "secret"
)

type Store struct {
	Name           string
	AllowedTenants []string
	Store          interface{}
	StoreType      string
}
