package log

type LoggerLevel string
type LoggerFormat string

const (
	InfoLevel  LoggerLevel = "info"
	ErrorLevel LoggerLevel = "error"
	DebugLevel LoggerLevel = "debug"
	WarnLevel  LoggerLevel = "warn"
	PanicLevel LoggerLevel = "panic"
)

const (
	TextFormat LoggerFormat = "text"
	JSONFormat LoggerFormat = "json"
)

type Config struct {
	Level  LoggerLevel
	Format LoggerFormat
}

func NewConfig(level LoggerLevel, format LoggerFormat) *Config {
	return &Config{
		Level:  level,
		Format: format,
	}
}
