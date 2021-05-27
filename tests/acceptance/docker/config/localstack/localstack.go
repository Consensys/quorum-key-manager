package localstack

const defaultLocalstackVaultImage = "localstack/localstack"
const DefaultHostPort = "4566"
const defaultHost = "localhost"
const defaultRegion = "eu-west-3"
const defaultAccessID = "test"
const defaultAccessKey = "test"

type Config struct {
	Image    string
	Port     string
	Host     string
	Region   string
	Services []string
}

func NewDefault() *Config {
	return &Config{
		Image:  defaultLocalstackVaultImage,
		Port:   DefaultHostPort,
		Host:   defaultHost,
		Region: defaultRegion,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (cfg *Config) SetRegion(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetServices(services []string) *Config {
	cfg.Services = services
	return cfg
}
