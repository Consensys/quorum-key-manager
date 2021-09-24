package utils

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	dockerhashicorp "github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/hashicorp"
)

const HashicorpPluginFilename = "quorum-hashicorp-vault-plugin"
const HashicorpPluginVersion = "v1.0.0"

func HashicorpContainer(logger log.Logger) (*dockerhashicorp.Config, error) {
	hashicorpHost := "localhost"
	hashicorpPort := strconv.Itoa(10000 + rand.Intn(10000))
	hashicorpToken := fmt.Sprintf("root_token_%v", strconv.Itoa(rand.Intn(10000)))

	pluginPath, err := getPluginPath()
	if err != nil {
		logger.WithError(err).Error("failed to get the current directory path")
		return nil, err
	}

	vaultContainer := dockerhashicorp.
		NewDefault().
		SetHostPort(hashicorpPort).
		SetRootToken(hashicorpToken).
		SetHost(hashicorpHost).
		SetPluginSourceDirectory(pluginPath)

	// Deal with darwin compliant plugin
	runtimeOS := runtime.GOOS
	switch runtimeOS {
	case "darwin":
		pluginPath += "/darwin"
		if _, err := os.Stat(pluginPath + "/" + HashicorpPluginFilename); os.IsNotExist(err) {
			logger.WithError(err).Error("cannot find required " + HashicorpPluginFilename + " file in " + pluginPath)
			return nil, err
		}
		vaultContainer.SetPluginSourceDirectory(pluginPath)
		logger.Info("using local Quorum Hashicorp Vault plugin", "path", pluginPath)

	default:
		logger.Info("downloading Quorum Hashicorp Vault plugin...", "path", pluginPath)
		pluginPath, err = vaultContainer.DownloadPlugin(HashicorpPluginFilename, HashicorpPluginVersion)
		if err != nil {
			logger.WithError(err).Error("cannot download hashicorp vault plugin")
			return nil, err
		}
		logger.Info("Quorum Hashicorp Vault plugin plugin downloaded", "path", pluginPath)
	}

	return vaultContainer, nil
}

func getPluginPath() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/plugins", currDir), nil
}
