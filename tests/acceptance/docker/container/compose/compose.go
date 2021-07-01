package compose

import (
	"context"
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	goreflect "reflect"
	"time"

	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config"
	hashConfig "github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config/hashicorp"
	hashVault "github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/container/hashicorp"
	"github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/container/reflect"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Compose struct {
	reflect *reflect.Reflect
}

func New(logger log.Logger) *Compose {
	factory := &Compose{
		reflect: reflect.New(),
	}

	factory.reflect.AddGenerator(goreflect.TypeOf(&hashConfig.Config{}), hashVault.New(logger))

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
