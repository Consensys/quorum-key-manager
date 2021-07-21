package flags

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(httpPortViperKey, httpPortDefault)
	_ = viper.BindEnv(httpPortViperKey, httpPortEnv)

	viper.SetDefault(healthPortViperKey, healthPortDefault)
	_ = viper.BindEnv(healthPortViperKey, healthPortEnv)

	viper.SetDefault(httpHostViperKey, httpHostDefault)
	_ = viper.BindEnv(httpHostViperKey, httpHostEnv)

	viper.SetDefault(httpsPortViperKey, httpsPortDefault)
	_ = viper.BindEnv(httpsPortViperKey, httpsPortEnv)

	viper.SetDefault(httpsHostViperKey, httpsHostDefault)
	_ = viper.BindEnv(httpsHostViperKey, httpsHostEnv)

	viper.SetDefault(tlsServerKeyViperKey, tlsServerKeyDefault)
	_ = viper.BindEnv(tlsServerKeyViperKey, tlsServerKeyEnv)

	viper.SetDefault(tlsServerCertViperKey, tlsServerCertDefault)
	_ = viper.BindEnv(tlsServerCertViperKey, tlsServerCertEnv)

}

const (
	tlsServerKeyFlag      = "tls-server-key"
	tlsServerKeyViperKey  = "tls.server.key"
	tlsServerKeyDefault   = "/tls/localhost.key"
	tlsServerKeyEnv       = "TLS_SERVER-KEY"
	tlsServerCertFlag     = "tls-server-cert"
	tlsServerCertViperKey = "tls.server.cert"
	tlsServerCertDefault  = "/tls/localhost.crt"
	tlsServerCertEnv      = "TLS_SERVER-CERT"
)

func tlsServerKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS key file location
Environment variable: %q`, tlsServerKeyEnv)
	f.String(tlsServerKeyFlag, tlsServerKeyDefault, desc)
	_ = viper.BindPFlag(tlsServerKeyViperKey, f.Lookup(tlsServerKeyFlag))
}

func tlsServerCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS cert file location
Environment variable: %q`, tlsServerCertEnv)
	f.String(tlsServerCertFlag, tlsServerCertDefault, desc)
	_ = viper.BindPFlag(tlsServerCertViperKey, f.Lookup(tlsServerCertFlag))
}

const (
	httpPortFlag      = "http-port"
	httpPortViperKey  = "http.port"
	httpPortDefault   = 8080
	httpPortEnv       = "HTTP_PORT"
	httpsPortFlag     = "https-port"
	httpsPortViperKey = "https.port"
	httpsPortDefault  = 8443
	httpsPortEnv      = "HTTPS_PORT"
)

func httpPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose API HTTP service
Environment variable: %q`, httpPortEnv)
	f.Uint32(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(httpPortViperKey, f.Lookup(httpPortFlag))
}

func httpsPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose API HTTPS service
Environment variable: %q`, httpsPortEnv)
	f.Uint32(httpsPortFlag, httpsPortDefault, desc)
	_ = viper.BindPFlag(httpsPortViperKey, f.Lookup(httpsPortFlag))
}

const (
	healthPortFlag     = "health-port"
	healthPortViperKey = "health.port"
	healthPortDefault  = 8081
	healthPortEnv      = "HEALTH_PORT"
)

func healthPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose Health HTTP service
Environment variable: %q`, healthPortEnv)
	f.Uint32(healthPortFlag, healthPortDefault, desc)
	_ = viper.BindPFlag(healthPortViperKey, f.Lookup(healthPortFlag))
}

const (
	httpHostFlag      = "http-host"
	httpHostViperKey  = "http.host"
	httpHostDefault   = ""
	httpHostEnv       = "HTTPS_HOST"
	httpsHostFlag     = "https-host"
	httpsHostViperKey = "https.host"
	httpsHostDefault  = ""
	httpsHostEnv      = "HTTPS_HOST"
)

// Hostname register a flag for HTTP server address
func httpHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Host to expose HTTP service
Environment variable: %q`, httpHostEnv)
	f.String(httpHostFlag, httpHostDefault, desc)
	_ = viper.BindPFlag(httpHostViperKey, f.Lookup(httpHostFlag))
}

// Hostname register a flag for HTTPS server address
func httpsHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Host to expose HTTPS service
Environment variable: %q`, httpsHostEnv)
	f.String(httpsHostFlag, httpsHostDefault, desc)
	_ = viper.BindPFlag(httpsHostViperKey, f.Lookup(httpHostFlag))
}

// Flags register flags for HashiCorp Hashicorp
func HTTPFlags(f *pflag.FlagSet) {
	httpHost(f)
	httpPort(f)
	healthPort(f)
	httpsHost(f)
	httpsPort(f)
	tlsServerCert(f)
	tlsServerKey(f)
}

func newHTTPConfig(vipr *viper.Viper) *server.Config {
	cfg := server.NewDefaultConfig()
	cfg.Port = vipr.GetUint32(httpPortViperKey)
	cfg.HealthzPort = vipr.GetUint32(healthPortViperKey)
	cfg.Host = vipr.GetString(httpHostViperKey)
	cfg.TLSPort = vipr.GetUint32(httpsPortViperKey)
	cfg.TLSHost = vipr.GetString(httpsHostViperKey)
	cfg.TLSCert = vipr.GetString(tlsServerCertViperKey)
	cfg.TLSKey = vipr.GetString(tlsServerKeyViperKey)

	return cfg
}
