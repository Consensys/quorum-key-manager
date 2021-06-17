package utils

import (
	dockerlocalstack "github.com/consensysquorum/quorum-key-manager/tests/acceptance/docker/config/localstack"
)

func LocalstackContainer() (*dockerlocalstack.Config, error) {
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
