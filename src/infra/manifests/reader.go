package manifests

import (
	"github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads manifests from filesystem
type Reader interface {
	Load() ([]*manifest.Manifest, error)
}
