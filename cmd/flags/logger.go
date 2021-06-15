package flags

import (
	"fmt"

	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(LevelViperKey, levelDefault)
	_ = viper.BindEnv(LevelViperKey, levelEnv)
	viper.SetDefault(FormatViperKey, formatDefault)
	_ = viper.BindEnv(FormatViperKey, formatEnv)
	viper.SetDefault(TimestampViperKey, timestampDefault)
	_ = viper.BindEnv(TimestampViperKey, timestampEnv)
}

const (
	levelFlag     = "log-old-level"
	LevelViperKey = "log-old.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

const (
	formatFlag     = "log-old-format"
	FormatViperKey = "log-old.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

const (
	timestampFlag     = "log-old-timestamp"
	TimestampViperKey = "log-old.timestamp"
	timestampDefault  = true
	timestampEnv      = "LOG_TIMESTAMP"
)

var ECSJsonFormatter = &logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyLevel: "log-old.level",
		logrus.FieldKeyMsg:   "message",
	},
}

func LoggerFlags(f *pflag.FlagSet) {
	level(f)
	format(f)
	timestamp(f)
}

// Level register flag for Level
func level(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, levelEnv)
	f.String(levelFlag, levelDefault, desc)
	_ = viper.BindPFlag(LevelViperKey, f.Lookup(levelFlag))
}

// Format register flag for Log Format
func format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = viper.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
}

func timestamp(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable logging with timestamp (only TEXT format).
Environment variable: %q`, timestampEnv)
	f.Bool(timestampFlag, timestampDefault, desc)
	_ = viper.BindPFlag(TimestampViperKey, f.Lookup(timestampFlag))
}

func newLoggerConfig(vipr *viper.Viper) *log_old.Config {
	return &log_old.Config{
		Level:     log_old.LoggerLevel(vipr.GetString(LevelViperKey)),
		Timestamp: vipr.GetBool(TimestampViperKey),
	}
}
