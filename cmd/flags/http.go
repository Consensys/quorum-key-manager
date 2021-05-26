package flags

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(httpPortViperKey, httpPortDefault)
	_ = viper.BindEnv(httpPortViperKey, httpPortEnv)

	viper.SetDefault(httpHostViperKey, httpHostDefault)
	_ = viper.BindEnv(httpHostViperKey, httpHostEnv)
}

const (
	httpPortFlag     = "http-port"
	httpPortViperKey = "server.port"
	httpPortDefault  = 8080
	httpPortEnv      = "HTTP_PORT"
)

func httpPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose HTTP service
Environment variable: %q`, httpPortEnv)
	f.Uint32(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(httpPortViperKey, f.Lookup(httpPortFlag))
}

const (
	httpHostFlag     = "http-host"
	httpHostViperKey = "server.host"
	httpHostDefault  = ""
	httpHostEnv      = "HTTP_HOST"
)

// Hostname register a flag for HTTP server address
func httpHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Host to expose HTTP service
Environment variable: %q`, httpHostEnv)
	f.String(httpHostFlag, httpHostDefault, desc)
	_ = viper.BindPFlag(httpHostViperKey, f.Lookup(httpHostFlag))
}

// Flags register flags for HashiCorp Hashicorp
func HTTPFlags(f *pflag.FlagSet) {
	httpHost(f)
	httpPort(f)
}

func newHTTPConfig(vipr *viper.Viper) *server.Config {
	cfg := server.NewDefaultConfig()
	cfg.Port = vipr.GetUint32(httpPortViperKey)
	cfg.Host = vipr.GetString(httpHostViperKey)
	return cfg
}
