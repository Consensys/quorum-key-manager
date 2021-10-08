package app

import manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"

type Config struct {
	Manifest *manifests.Config
}
