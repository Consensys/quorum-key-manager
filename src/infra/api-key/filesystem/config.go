package filesystem

type Config struct {
	Path string
}

func NewConfig(path string) *Config {
	return &Config{
		Path: path,
	}
}
