package container

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type DockerContainerFactory interface {
	GenerateContainerConfig(ctx context.Context, configuration interface{}) (*container.Config, *container.HostConfig, *network.NetworkingConfig, error)
	WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error
}
