package flags

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oicd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authOICDCACertFileViperKey, authOICDCACertFileEnv)
	_ = viper.BindEnv(AuthOICDCAKeyFileViperKey, authOICDCAKeyFileEnv)

	viper.SetDefault(authOICDClaimUsernameViperKey, authOICDClaimUsernameDefault)
	_ = viper.BindEnv(authOICDClaimUsernameViperKey, authOICDClaimUsernameEnv)

	viper.SetDefault(authOICDClaimGroupViperKey, authOICDClaimGroupDefault)
	_ = viper.BindEnv(authOICDClaimGroupViperKey, authOICDClaimGroupEnv)

}

const (
	authOICDCACertFileFlag     = "auth-oidc-ca-cert"
	authOICDCACertFileViperKey = "auth.oidc.ca.cert"
	authOICDCACertFileDefault  = ""
	authOICDCACertFileEnv      = "AUTH_OICD_CA_CERT"
)

const (
	authOICDCAKeyFileFlag     = "auth-oidc-ca-key"
	AuthOICDCAKeyFileViperKey = "auth.oidc.ca.key"
	authOICDCAKeyFileDefault  = ""
	authOICDCAKeyFileEnv      = "AUTH_OICD_CA_KEY"
)

const (
	authOICDClaimUsernameFlag     = "auth-oidc-claim-username"
	authOICDClaimUsernameViperKey = "auth.oidc.claim.username"
	authOICDClaimUsernameDefault  = "qkm.auth.username"
	authOICDClaimUsernameEnv      = "AUTH_OICD_CLAIM_USERNAME"
)

const (
	authOICDClaimGroupFlag     = "auth-oidc-claim-groups"
	authOICDClaimGroupViperKey = "auth.oidc.claim.groups"
	authOICDClaimGroupDefault  = "qkm.auth.groups"
	authOICDClaimGroupEnv      = "AUTH_OICD_CLAIM_GROUPS"
)

func AuthFlags(f *pflag.FlagSet) {
	authOICDCAFile(f)
	authOICDClaimUsername(f)
	authOICDClaimGroups(f)
}

// Use only on generate-token utils 
func AuthOICDCertKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authorization OICD CA Cert filepath.
Environment variable: %q`, authOICDCAKeyFileEnv)
	f.String(authOICDCAKeyFileFlag, authOICDCAKeyFileDefault, desc)
	_ = viper.BindPFlag(AuthOICDCAKeyFileViperKey, f.Lookup(authOICDCAKeyFileFlag))
}

func authOICDCAFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authorization OICD CA Key filepath.
Environment variable: %q`, authOICDClaimUsernameEnv)
	f.String(authOICDClaimUsernameFlag, authOICDClaimUsernameDefault, desc)
	_ = viper.BindPFlag(authOICDClaimUsernameViperKey, f.Lookup(authOICDClaimUsernameFlag))
}

func authOICDClaimUsername(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authorization OICD username claims.
Environment variable: %q`, authOICDClaimGroupEnv)
	f.String(authOICDClaimGroupFlag, authOICDClaimGroupDefault, desc)
	_ = viper.BindPFlag(authOICDClaimGroupViperKey, f.Lookup(authOICDClaimGroupFlag))
}

func authOICDClaimGroups(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Authorization OICD group claims.
Environment variable: %q`, authOICDCACertFileEnv)
	f.String(authOICDCACertFileFlag, authOICDCACertFileDefault, desc)
	_ = viper.BindPFlag(authOICDCACertFileViperKey, f.Lookup(authOICDCACertFileFlag))
}

func NewAuthConfig(vipr *viper.Viper) (*auth.Config, error) {
	caFile := vipr.GetString(authOICDCACertFileViperKey)
	_, err := os.Stat(caFile)
	
	var oicdCfg = &oicd.Config{}
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		caFileContent, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		oicdCfg = oicd.NewConfig(string(caFileContent), vipr.GetString(authOICDClaimUsernameViperKey), 
			vipr.GetString(authOICDClaimGroupViperKey))
	}

	return &auth.Config{OICD: oicdCfg}, nil
}
