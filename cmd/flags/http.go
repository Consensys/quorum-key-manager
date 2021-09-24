package flags

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

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
	enableHTTPSFlag     = "https-enable"
	enableHTTPSViperKey = "https-enable"
	enableHTTPSDefault  = false
	enableHTTPSEnv      = "HTTPS_ENABLED"
)

// Hostname register a flag for HTTP server address
func enableHTTPS(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable https server
Environment variable: %q`, enableHTTPSEnv)
	f.Bool(enableHTTPSFlag, enableHTTPSDefault, desc)
	_ = viper.BindPFlag(enableHTTPSViperKey, f.Lookup(enableHTTPSFlag))
}

const (
	httpServerKeyFlag     = "https-server-key"
	httpServerKeyViperKey = "https.server.key"
	httpServerKeyDefault  = ""
	httpServerKeyEnv      = "HTTPS_SERVER_KEY"
)

func httpServerKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS key file location
Environment variable: %q`, httpServerKeyEnv)
	f.String(httpServerKeyFlag, httpServerKeyDefault, desc)
	_ = viper.BindPFlag(httpServerKeyViperKey, f.Lookup(httpServerKeyFlag))
}

const (
	httpServerCertFlag     = "https-server-cert"
	httpServerCertViperKey = "https.server.cert"
	httpServerCertDefault  = ""
	httpServerCertEnv      = "HTTPS_SERVER_CERT"
)

func httpServerCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS cert file location
Environment variable: %q`, httpServerCertEnv)
	f.String(httpServerCertFlag, httpServerCertDefault, desc)
	_ = viper.BindPFlag(httpServerCertViperKey, f.Lookup(httpServerCertFlag))
}

// HTTPFlags register flags for HTTPS server
func HTTPFlags(f *pflag.FlagSet) {
	httpHost(f)
	httpPort(f)
	healthPort(f)
	enableHTTPS(f)
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
			MinVersion:               tls.VersionTLS13,
			ClientAuth:               tls.VerifyClientCertIfGiven,
			PreferServerCipherSuites: true,
		}

		certFile := vipr.GetString(httpServerCertViperKey)
		keyFile := vipr.GetString(httpServerKeyViperKey)

		var err error
		cfg.TLSConfig.Certificates = make([]tls.Certificate, 1)
		cfg.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read http ssl certificates '%s' or '%s': %w", certFile, keyFile, err)
		}

		clientCAPool, err := clientCAPool(vipr)
		if err != nil {
			return nil, err
		} else if clientCAPool != nil {
			cfg.TLSConfig.ClientCAs = clientCAPool
		}
	}

	return cfg, nil
}

func clientCAPool(vipr *viper.Viper) (*x509.CertPool, error) {
	caFile := vipr.GetString(authTLSCertsFileViperKey)
	if caFile == "" {
		return nil, nil
	}
	_, err := os.Stat(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca file. %s", err.Error())
	}

	caFileContent, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caFileContent)

	return caCertPool, nil
}
