package flags

import (
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(AuthOIDCPrivKeyViperKey, authOIDCPrivKeyEnv)
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)
	_ = viper.BindEnv(AuthOIDCAudienceViperKey, authOIDCAudienceEnv)
}

const (
	authOIDCIssuerURLFlag     = "auth-oidc-issuer-url"
	authOIDCIssuerURLViperKey = "auth.oidc.issuer.url"
	authOIDCIssuerURLDefault  = ""
	authOIDCIssuerURLEnv      = "AUTH_OIDC_ISSUER_URL"
)

const (
	AuthOIDCPrivKeyViperKey = "auth.oidc.priv.key"
	authOIDCPrivKeyEnv      = "AUTH_OIDC_PRIV_KEY"
)

const (
	authOIDCAudienceFlag     = "auth-oidc-audience"
	AuthOIDCAudienceViperKey = "auth.oidc.audience"
	authOIDCAudienceEnv      = "AUTH_OIDC_AUDIENCE"
)

func OIDCFlags(f *pflag.FlagSet) {
	authOIDCIssuerServer(f)
	authOIDCAudience(f)
}

func authOIDCIssuerServer(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect issuer server domain (ie. https://quorum-key-manager.eu.auth0.com).
Environment variable: %q`, authOIDCIssuerURLEnv)
	f.String(authOIDCIssuerURLFlag, authOIDCIssuerURLDefault, desc)
	_ = viper.BindPFlag(authOIDCIssuerURLViperKey, f.Lookup(authOIDCIssuerURLFlag))
}

func authOIDCAudience(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Expected audience ("aud" field) of JWT tokens.
Environment variable: %q`, authOIDCAudienceEnv)
	f.StringArray(authOIDCAudienceFlag, []string{}, desc)
	_ = viper.BindPFlag(AuthOIDCAudienceViperKey, f.Lookup(authOIDCAudienceFlag))
}

func NewOIDCConfig(vipr *viper.Viper) *jose.Config {
	issuerURL := vipr.GetString(authOIDCIssuerURLViperKey)

	if issuerURL != "" {
		return jose.NewConfig(
			vipr.GetString(authOIDCIssuerURLViperKey),
			vipr.GetStringSlice(AuthOIDCAudienceViperKey),
			5*time.Minute, // TODO: Make the cache ttl an ENV var if needed
		)
	}

	return nil
}
