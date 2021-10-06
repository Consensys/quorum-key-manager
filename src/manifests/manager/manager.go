package manager

import (
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager manages manifests
type Manager interface {
	Load() ([]manifest.Message, error)
}
