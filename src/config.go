package src

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/http"
)

type Config struct {
	Logger    *log.Config
	HTTP      *http.Config
	Manifests []*manifest.Manifest
}
