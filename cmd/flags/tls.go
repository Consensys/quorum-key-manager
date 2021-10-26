package flags

import (
	"crypto/x509"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/auth"
	authtls "github.com/consensys/quorum-key-manager/src/auth/service/authenticator/tls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

func init() {
	_ = viper.BindEnv(authTLSCertsFileViperKey, authTLSCertsFileEnv)
}

const (
	authTLSCertsFileFlag     = "auth-tls-ca"
	authTLSCertsFileViperKey = "auth.tls.ca"
	authTLSCertsFileDefault  = ""
	authTLSCertsFileEnv      = "AUTH_TLS_CA"
)

func TLSFlags(f *pflag.FlagSet) {
	authTLSCertFile(f)
}

func authTLSCertFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Authenticator Cert filepath.
Environment variable: %q`, authTLSCertsFileEnv)
	f.String(authTLSCertsFileFlag, authTLSCertsFileDefault, desc)
	_ = viper.BindPFlag(authTLSCertsFileViperKey, f.Lookup(authTLSCertsFileFlag))
}

func NewTLSConfig(vipr *viper.Viper) (*auth.Config, error) {
	var tlsCfg *authtls.Config
	tlsAuthCAs, err := tlsAuthCerts(vipr)
	if err != nil {
		return nil, err
	}

	tlsCfg = authtls.NewConfig(tlsAuthCAs)

	return &auth.Config{
		OIDC:     oidcCfg,
		APIKEY:   apiKeyCfg,
		TLS:      tlsCfg,
		Manifest: NewManifestConfig(vipr),
	}, nil
}

func tlsAuthCerts(vipr *viper.Viper) (*x509.CertPool, error) {
	caFile := vipr.GetString(authTLSCertsFileViperKey)
	if caFile == "" {
		return nil, nil
	}

	_, err := os.Stat(caFile)
	if err != nil {
		return nil, err
	}

	caFileContent, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caFileContent)
	if !ok {
		return nil, fmt.Errorf("failed to append cert to pool")
	}

	return caCertPool, nil
}
