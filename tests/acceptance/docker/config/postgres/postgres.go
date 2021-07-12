package postgres

const (
	DefaultPostgresImage = "postgres:13.3-alpine"
	defaultPassword      = "postgres"
	defaultHostPort      = "5432"
)

type Config struct {
	Image    string
	Port     string
	Password string
}

func NewDefault() *Config {
	return &Config{
		Image:    DefaultPostgresImage,
		Port:     defaultHostPort,
		Password: defaultPassword,
	}
}

func (c *Config) SetPort(port string) *Config {
	c.Port = port
	return c
}
