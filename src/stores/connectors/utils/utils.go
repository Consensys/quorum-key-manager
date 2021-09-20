package utils

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
)

type Connector struct {
	logger log.Logger
}

var _ stores.Utilities = Connector{}

func NewConnector(logger log.Logger) *Connector {
	return &Connector{
		logger: logger,
	}
}
