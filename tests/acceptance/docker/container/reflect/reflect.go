package reflect

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/consensys/quorum-key-manager/tests/acceptance/docker/container"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Reflect struct {
	generators map[reflect.Type]container.DockerContainerFactory
}

func New() *Reflect {
	return &Reflect{
		generators: make(map[reflect.Type]container.DockerContainerFactory),
	}
}

func (gen *Reflect) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	generator, ok := gen.generators[reflect.TypeOf(configuration)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("no container config generator for configuration of type %T (consider adding one)", configuration)
	}

	return generator.GenerateContainerConfig(ctx, configuration)
}

func (gen *Reflect) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	generator, ok := gen.generators[reflect.TypeOf(configuration)]
	if !ok {
		return fmt.Errorf("no container config generator for configuration of type %T (consider adding one)", configuration)
	}

	return generator.WaitForService(ctx, configuration, timeout)
}

func (gen *Reflect) AddGenerator(typ reflect.Type, generator container.DockerContainerFactory) {
	gen.generators[typ] = generator
}
