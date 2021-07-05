package utils

import (
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/log"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	dockerhashicorp "github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/hashicorp"
)

const HashicorpPluginFilename = "orchestrate-hashicorp-vault-plugin"
const HashicorpPluginVersion = "v0.0.11-alpha.3"

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

	//Deal with darwin compliant plugin
	runtime := runtime.GOOS
	switch runtime {
	case "darwin":
		pluginPath += "/darwin"
		if _, err := os.Stat(pluginPath + "/" + HashicorpPluginFilename); os.IsNotExist(err) {
			logger.WithError(err).Error("cannot find required " + HashicorpPluginFilename + " file in " + pluginPath)
			return nil, err
		}
		vaultContainer.SetPluginSourceDirectory(pluginPath)
		logger.Info("using local orchestrate plugin", "path", pluginPath)

	default:
		pluginPath, err = vaultContainer.DownloadPlugin(HashicorpPluginFilename, HashicorpPluginVersion)
		if err != nil {
			logger.WithError(err).Error("cannot download hashicorp vault plugin")
			return nil, err
		}
		logger.Info("orchestrate plugin downloaded", "path", pluginPath)
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
