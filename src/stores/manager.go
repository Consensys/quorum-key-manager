package stores

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows managing multiple stores
type Manager interface {
	Stores() Stores
}
