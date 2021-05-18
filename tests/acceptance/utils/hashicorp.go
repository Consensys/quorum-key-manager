package utils

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	dockerhashicorp "github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/docker/container/hashicorp"
)

const HashicorpPluginFilename = "orchestrate-hashicorp-vault-plugin"
const HashicorpPluginVersion = "v0.0.11-alpha.1"

func HashicorpContainer(ctx context.Context) (*dockerhashicorp.Config, error) {
	logger := log.FromContext(ctx)

	hashicorpHost := "localhost"
	hashicorpPort := strconv.Itoa(10000 + rand.Intn(10000))
	hashicorpToken := fmt.Sprintf("root_token_%v", strconv.Itoa(rand.Intn(10000)))

	pluginPath, err := getPluginPath(logger)
	if err != nil {
		return nil, err
	}

	vaultContainer := dockerhashicorp.
		NewDefault().
		SetHostPort(hashicorpPort).
		SetRootToken(hashicorpToken).
		SetHost(hashicorpHost).
		SetPluginSourceDirectory(pluginPath)

	pluginPath, err = vaultContainer.DownloadPlugin(HashicorpPluginFilename, HashicorpPluginVersion)
	if err != nil {
		logger.WithError(err).Error("cannot download hashicorp vault plugin")
		return nil, err
	}
	logger.WithField("path", pluginPath).Info("orchestrate plugin downloaded")

	return vaultContainer, nil
}

func getPluginPath(logger *log.Logger) (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		logger.WithError(err).Error("failed to get the current directory path")
		return "", err
	}

	return fmt.Sprintf("%s/plugins", currDir), nil
}
