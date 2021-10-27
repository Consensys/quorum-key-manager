package apikey

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads manifests from filesystem
type Reader interface {
	Load(ctx context.Context) (map[string]*entities.UserClaims, error)
}
