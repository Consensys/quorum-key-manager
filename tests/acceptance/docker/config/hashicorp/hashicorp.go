package hashicorp

const defaultHashicorpVaultImage = "consensys/quorum-hashicorp-vault-plugin:v1.1.4"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"
const defaultPluginMountPath = "quorum"

type Config struct {
	Image           string
	Host            string
	Port            string
	RootToken       string
	PluginMountPath string
}

func NewDefault() *Config {
	return &Config{
		Image:           defaultHashicorpVaultImage,
		Port:            defaultHostPort,
		RootToken:       defaultRootToken,
		Host:            defaultHost,
		PluginMountPath: defaultPluginMountPath,
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

func (cfg *Config) SetMountPath(mountPath string) *Config {
	cfg.PluginMountPath = mountPath
	return cfg
}

