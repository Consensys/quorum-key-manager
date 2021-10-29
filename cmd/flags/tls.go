package flags

import (
	"fmt"

	tls "github.com/consensys/quorum-key-manager/src/infra/tls/filesystem"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

func NewTLSConfig(vipr *viper.Viper) *tls.Config {
	path := vipr.GetString(authTLSCertsFileViperKey)

	if path != "" {
		return tls.NewConfig(vipr.GetString(authTLSCertsFileViperKey))
	}

	return nil
}
