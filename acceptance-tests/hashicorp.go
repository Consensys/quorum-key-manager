package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	dockerhashicorp "github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/container/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

const hashicorpPluginFilename = "orchestrate-hashicorp-vault-plugin"
const hashicorpPluginVersion = "v0.0.10-alpha.5"

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

	pluginPath, err = vaultContainer.DownloadPlugin(hashicorpPluginFilename, hashicorpPluginVersion)
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
