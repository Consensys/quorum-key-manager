package compose

import (
	"context"
	"fmt"
	goreflect "reflect"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker/container/hashicorp"

	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker/config"
	"github.com/ConsenSysQuorum/quorum-key-manager/integration-tests/docker/container/reflect"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Compose struct {
	reflect *reflect.Reflect
}

func New() *Compose {
	factory := &Compose{
		reflect: reflect.New(),
	}

	factory.reflect.AddGenerator(goreflect.TypeOf(&hashicorp.Config{}), &hashicorp.Vault{})

	return factory
}

func (gen *Compose) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*config.Container)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, nil, nil, err
	}

	return gen.reflect.GenerateContainerConfig(ctx, field)
}

func (gen *Compose) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*config.Container)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return err
	}

	return gen.reflect.WaitForService(ctx, field, timeout)
}
