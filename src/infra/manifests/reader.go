package manifests

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads manifests
type Reader interface {
	Load(ctx context.Context) (map[string][]entities.Manifest, error)
}
