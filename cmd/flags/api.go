package flags

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(apiPortViperKey, apiPortDefault)
	_ = viper.BindEnv(apiPortViperKey, apiPortEnv)

	viper.SetDefault(apiHostViperKey, apiHostDefault)
	_ = viper.BindEnv(apiHostViperKey, apiHostEnv)
}

const (
	apiPortFlag     = "api-port"
	apiPortViperKey = "api.port"
	apiPortDefault  = 8080
	apiPortEnv      = "API_PORT"
)

func apiPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose API service
Environment variable: %q`, apiPortEnv)
	f.Uint32(apiPortFlag, apiPortDefault, desc)
	_ = viper.BindPFlag(apiPortViperKey, f.Lookup(apiPortFlag))
}

const (
	apiHostFlag     = "api-host"
	apiHostViperKey = "api.host"
	apiHostDefault  = "localhost"
	apiHostEnv      = "API_HOST"
)

// Hostname register a flag for HTTP server address
func apiHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Host to expose API service
Environment variable: %q`, apiHostEnv)
	f.String(apiHostFlag, apiHostDefault, desc)
	_ = viper.BindPFlag(apiHostViperKey, f.Lookup(apiHostFlag))
}

// Flags register flags for HashiCorp Hashicorp
func APIFlags(f *pflag.FlagSet) {
	apiHost(f)
	apiPort(f)
}

func newAPIConfig(vipr *viper.Viper) *api.Config {
	cfg := api.NewDefaultConfig()
	cfg.Port = vipr.GetUint32(apiPortViperKey)
	cfg.Host = vipr.GetString(apiHostViperKey)
	return cfg
}
