package client

type Config struct {
	URL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}
