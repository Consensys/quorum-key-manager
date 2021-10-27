package manifests

import (
	"context"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads manifests
type Reader interface {
	Load(ctx context.Context) ([]*manifest.Manifest, error)
}
