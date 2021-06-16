package flags

import (
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"

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
	viper.SetDefault(logModeViperKey, logModeDefault)
	_ = viper.BindEnv(logModeViperKey, logModeEnv)
}

const (
	levelFlag     = "log-level"
	LevelViperKey = "log.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

const (
	formatFlag     = "log-format"
	FormatViperKey = "log.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

const (
	timestampFlag     = "log-timestamp"
	TimestampViperKey = "log.timestamp"
	timestampDefault  = true
	timestampEnv      = "LOG_TIMESTAMP"
)

const (
	logModeFlag     = "log-mode"
	logModeViperKey = "log.timestamp"
	logModeDefault  = "production"
	logModeEnv      = "LOG_MODE"
)

func LoggerFlags(f *pflag.FlagSet) {
	level(f)
	format(f)
	timestamp(f)
	mode(f)
}

// Level register flag for Level
func level(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "error", "warn", "info", "debug"}, levelEnv)
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

func mode(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log mode (one of %q).
Environment variable: %q`, []log.LoggerMode{log.DevelopmentMode, log.ProductionMode}, logModeEnv)
	f.String(logModeFlag, logModeDefault, desc)
	_ = viper.BindPFlag(logModeViperKey, f.Lookup(logModeFlag))
}

func newLoggerConfig(vipr *viper.Viper) *log.Config {
	return log.NewConfig(
		log.LoggerLevel(vipr.GetString(LevelViperKey)),
		vipr.GetBool(TimestampViperKey),
		log.LoggerMode(vipr.GetString(logModeViperKey)),
	)
}
