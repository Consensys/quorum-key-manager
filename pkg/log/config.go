package log

type LoggerLevel string
type LoggerMode string

const (
	InfoLevel  LoggerLevel = "info"
	ErrorLevel LoggerLevel = "error"
	DebugLevel LoggerLevel = "debug"
	WarnLevel  LoggerLevel = "warn"
	PanicLevel LoggerLevel = "panic"
)

const (
	DevelopmentMode LoggerMode = "development"
	ProductionMode  LoggerMode = "production"
)

type Config struct {
	Level     LoggerLevel
	Timestamp bool
	Mode      LoggerMode
}

func NewConfig(level LoggerLevel, timestamp bool, mode LoggerMode) *Config {
	return &Config{
		Level:     level,
		Timestamp: timestamp,
		Mode:      mode,
	}
}
