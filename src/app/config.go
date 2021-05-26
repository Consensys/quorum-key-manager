package src

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

type Config struct {
	Logger       *log.Config
	HTTP         *server.Config
	ManifestPath string
}
