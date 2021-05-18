package utils

import (
	"context"

	dockerlocalstack "github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/docker/container/localstack"
)

func LocalstackContainer(ctx context.Context) (*dockerlocalstack.Config, error) {

	localstackHost := "localhost"
	localstackPort := "4566"
	localstackServices := []string{"s3", "kms", "secretsmanager"}

	vaultContainer := dockerlocalstack.
		NewDefault().
		SetHostPort(localstackPort).
		SetHost(localstackHost).
		SetServices(localstackServices)

	return vaultContainer, nil
}
