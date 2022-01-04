package hashicorp

const defaultHashicorpVaultImage = "consensys/quorum-hashicorp-vault-plugin:v1.1.4"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"
const defaultMountPath = "orchestrate"

type Config struct {
	Image     string
	Host      string
	Port      string
	RootToken string
	MonthPath string
}

func NewDefault() *Config {
	return &Config{
		Image:     defaultHashicorpVaultImage,
		Port:      defaultHostPort,
		RootToken: defaultRootToken,
		Host:      defaultHost,
		MonthPath: defaultMountPath,
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
	cfg.MonthPath = mountPath
	return cfg
}

