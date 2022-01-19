package flags

import (
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)
	_ = viper.BindEnv(AuthOIDCAudienceViperKey, authOIDCAudienceEnv)
	_ = viper.BindEnv(authOIDCQKMClaimsViperKey, authOIDCQKMClaimsEnv)
}

const (
	authOIDCIssuerURLFlag     = "auth-oidc-issuer-url"
	authOIDCIssuerURLViperKey = "auth.oidc.issuer.url"
	authOIDCIssuerURLDefault  = ""
	authOIDCIssuerURLEnv      = "AUTH_OIDC_ISSUER_URL"
)

const (
	authOIDCAudienceFlag     = "auth-oidc-audience"
	AuthOIDCAudienceViperKey = "auth.oidc.audience"
	authOIDCAudienceEnv      = "AUTH_OIDC_AUDIENCE"
)

const (
	authOIDCQKMClaimsFlag     = "auth-oidc-qkm-claim"
	authOIDCQKMClaimsViperKey = "auth.oidc.qkm.claims"
	authOIDCQKMClaimsEnv      = "AUTH_OIDC_QKM_CLAIMS"
)

func OIDCFlags(f *pflag.FlagSet) {
	authOIDCIssuerServer(f)
	authOIDCAudience(f)
	authOIDCQKMClaimsPath(f)
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
	f.StringSlice(authOIDCAudienceFlag, []string{}, desc)
	_ = viper.BindPFlag(AuthOIDCAudienceViperKey, f.Lookup(authOIDCAudienceFlag))
}

func authOIDCQKMClaimsPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to for Quorum Key Manager custom claims in the JWT.
Environment variable: %q`, authOIDCQKMClaimsEnv)
	f.String(authOIDCQKMClaimsFlag, "", desc)
	_ = viper.BindPFlag(authOIDCQKMClaimsViperKey, f.Lookup(authOIDCQKMClaimsFlag))
}

func NewOIDCConfig(vipr *viper.Viper) *jose.Config {
	issuerURL := vipr.GetString(authOIDCIssuerURLViperKey)

	if issuerURL != "" {
		return jose.NewConfig(
			issuerURL,
			vipr.GetStringSlice(AuthOIDCAudienceViperKey),
			5*time.Minute, // TODO: Make the cache ttl an ENV var if needed
		)
	}

	return nil
}
