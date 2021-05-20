package hashicorp

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultHashicorpVaultImage = "library/vault:1.6.2"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"

type Config struct {
	Image                 string
	Port                  string
	RootToken             string
	Host                  string
	PluginSourceDirectory string
}

func NewDefault() *Config {
	return &Config{
		Image:     defaultHashicorpVaultImage,
		Port:      defaultHostPort,
		RootToken: defaultRootToken,
		Host:      defaultHost,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetRootToken(rootToken string) *Config {
	cfg.RootToken = rootToken
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (cfg *Config) SetPluginSourceDirectory(dir string) *Config {
	cfg.PluginSourceDirectory = dir
	return cfg
}

func (cfg *Config) DownloadPlugin(filename, version string) (string, error) {
	url := fmt.Sprintf("https://github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/releases/download/%s/orchestrate-hashicorp-vault-plugin", version)

	pluginPath := fmt.Sprintf("%s/%s", cfg.PluginSourceDirectory, filename)
	err := downloadPlugin(pluginPath, url)
	if err != nil {
		return "", err
	}
	return pluginPath, nil
}

func downloadPlugin(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = os.Chmod(filepath, 0777)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
