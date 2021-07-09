package postgres

//go:generate mockgen -source=postgres.go -destination=mocks/postgres.go -package=mocks

type Client interface {
	Insert(model ...interface{}) error
	SelectPK(model ...interface{}) error
	SelectDeletedPK(model ...interface{}) error
	Select(model ...interface{}) error
	SelectDeleted(model ...interface{}) error
	UpdatePK(model ...interface{}) error
	DeletePK(model ...interface{}) error
	ForceDeletePK(model ...interface{}) error
}
