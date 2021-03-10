package src

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/logger"
)

type Config struct {
	Logger *logger.Config
	HTTP   *api.Config
}
