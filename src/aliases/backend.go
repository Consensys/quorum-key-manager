package aliases

//go:generate mockgen -destination=mock/backend.go -package=mock . Backend

type Backend interface {
	Repository
	Parser
}
