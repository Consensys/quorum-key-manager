package entities

type ETH1Account struct {
	ID       string
	Address  string
	Metadata *Metadata
	Tags     map[string]string
}
