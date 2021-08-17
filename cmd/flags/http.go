package flags

import (
	"crypto/tls"
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

	viper.SetDefault(enableHTTPSViperKey, enableHTTPSDefault)
	_ = viper.BindEnv(enableHTTPSViperKey, enableHTTPSEnv)

	viper.SetDefault(httpServerKeyViperKey, httpServerKeyDefault)
	_ = viper.BindEnv(httpServerKeyViperKey, httpServerKeyEnv)

	viper.SetDefault(httpServerCertViperKey, httpServerCertDefault)
	_ = viper.BindEnv(httpServerCertViperKey, httpServerCertEnv)

}

const (
	httpPortFlag     = "http-port"
	httpPortViperKey = "http.port"
	httpPortDefault  = 8080
	httpPortEnv      = "HTTP_PORT"
)

func httpPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose API HTTP service
Environment variable: %q`, httpPortEnv)
	f.Uint32(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(httpPortViperKey, f.Lookup(httpPortFlag))
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
	httpHostFlag     = "http-host"
	httpHostViperKey = "http.host"
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

const (
	enableHTTPSFlag     = "enable-https"
	enableHTTPSViperKey = "enable.https"
	enableHTTPSDefault  = false
	enableHTTPSEnv      = "HTTP_SERVER_SSL"
)

// Hostname register a flag for HTTP server address
func enableHTTPs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable https server
Environment variable: %q`, enableHTTPSEnv)
	f.Bool(enableHTTPSFlag, enableHTTPSDefault, desc)
	_ = viper.BindPFlag(enableHTTPSViperKey, f.Lookup(enableHTTPSFlag))
}

const (
	httpServerKeyFlag     = "tls-server-key"
	httpServerKeyViperKey = "tls.server.key"
	httpServerKeyDefault  = ""
	httpServerKeyEnv      = "HTTP_SERVER_KEY"
)

func httpServerKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS key file location
Environment variable: %q`, httpServerKeyEnv)
	f.String(httpServerKeyFlag, httpServerKeyDefault, desc)
	_ = viper.BindPFlag(httpServerKeyViperKey, f.Lookup(httpServerKeyFlag))
}

const (
	httpServerCertFlag     = "tls-server-cert"
	httpServerCertViperKey = "tls.server.cert"
	httpServerCertDefault  = ""
	httpServerCertEnv      = "HTTP_SERVER_CERT"
)

func httpServerCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS cert file location
Environment variable: %q`, httpServerCertEnv)
	f.String(httpServerCertFlag, httpServerCertDefault, desc)
	_ = viper.BindPFlag(httpServerCertViperKey, f.Lookup(httpServerCertFlag))
}

// Flags register flags for HashiCorp Hashicorp
func HTTPFlags(f *pflag.FlagSet) {
	httpHost(f)
	httpPort(f)
	healthPort(f)
	enableHTTPs(f)
	httpServerCert(f)
	httpServerKey(f)
}

func newHTTPConfig(vipr *viper.Viper) (*server.Config, error) {
	cfg := server.NewDefaultConfig()
	cfg.Port = vipr.GetUint32(httpPortViperKey)
	cfg.HealthzPort = vipr.GetUint32(healthPortViperKey)
	cfg.Host = vipr.GetString(httpHostViperKey)

	isSSL := vipr.GetBool(enableHTTPSViperKey)
	if isSSL {
		cfg.TLSConfig = &tls.Config{
			ClientAuth: tls.VerifyClientCertIfGiven,
		}
		certFile := vipr.GetString(httpServerCertViperKey)
		keyFile := vipr.GetString(httpServerKeyViperKey)
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read http ssl certificates '%s' or '%s': %w", certFile, keyFile, err)
		}
		
		clientCAPool, err := clientCAPool(vipr)
		if clientCAPool != nil {
			cfg.TLSConfig.ClientCAs = clientCAPool
		}
		
		cfg.TLSConfig.Certificates = append(cfg.TLSConfig.Certificates, cert)
	}

	return cfg, nil
}
