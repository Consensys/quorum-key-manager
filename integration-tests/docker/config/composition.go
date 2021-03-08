package config

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/common/utils"
	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker/container/hashicorp"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	HashicorpVault *hashicorp.Config
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}
