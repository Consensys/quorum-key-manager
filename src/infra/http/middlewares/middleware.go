package middlewares

import "net/http"

//go:generate mockgen -source=middleware.go -destination=mock/middleware.go -package=mock

type Middleware interface {
	Handler(next http.Handler) http.Handler
}
