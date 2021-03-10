package src

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
)

type Config struct {
	Logger *log.Config
	HTTP   *api.Config
}
