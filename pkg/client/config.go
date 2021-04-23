package client

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	urlFlag     = "key-manager-url"
	URLViperKey = "key.manager.url"
	urlDefault  = "http://localhost:8081"
	urlEnv      = "KEY_MANAGER_URL"
)

const (
	MetricsURLViperKey = "key.manager.metrics.url"
	metricsURLDefault  = "http://localhost:8082"
	metricsURLEnv      = "KEY_MANAGER_METRICS_URL"
)

func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Key Manager HTTP endpoint. 
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

func Flags(f *pflag.FlagSet) {
	URL(f)
}

type Config struct {
	URL        string
	MetricsURL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL:        url,
		MetricsURL: metricsURLDefault,
	}
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL:        vipr.GetString(URLViperKey),
		MetricsURL: vipr.GetString(MetricsURLViperKey),
	}
}
