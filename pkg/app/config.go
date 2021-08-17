package app

import (
	"github.com/consensys/quorum-key-manager/pkg/http/server"
)

type Config struct {
	HTTP  *server.Config
}
