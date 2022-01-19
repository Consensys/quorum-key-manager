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
	authOIDCQKMClaimsFlag     = "auth-oidc-qkm-claims"
	authOIDCQKMClaimsViperKey = "auth.oidc.qkm.claims"
	authOIDCQKMClaimsEnv      = "AUTH_OIDC_QKM_CLAIMS"
)

const (
	authOIDCPermissionsClaimsFlag     = "auth-oidc-permissions-claims"
	authOIDCPermissionsClaimsViperKey = "auth.oidc.permissions.claims"
	authOIDCPermissionsClaimsEnv      = "AUTH_OIDC_PERMISSIONS_CLAIMS"
)

const (
	authOIDCRolesClaimsFlag     = "auth-oidc-roles-claim"
	authOIDCRolesClaimsViperKey = "auth.oidc.roles.claims"
	authOIDCRolesClaimsEnv      = "AUTH_OIDC_ROLES_CLAIMS"
)

func OIDCFlags(f *pflag.FlagSet) {
	authOIDCIssuerServer(f)
	authOIDCAudience(f)
	authOIDCQKMClaimsPath(f)
	authOIDCPermissionsClaimsPath(f)
	authOIDCRolesClaimsPath(f)
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

func authOIDCPermissionsClaimsPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to for permissions claims in the JWT (default 'scope').
Environment variable: %q`, authOIDCPermissionsClaimsEnv)
	f.String(authOIDCPermissionsClaimsFlag, "", desc)
	_ = viper.BindPFlag(authOIDCPermissionsClaimsViperKey, f.Lookup(authOIDCPermissionsClaimsFlag))
}

func authOIDCRolesClaimsPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to for roles claims in the JWT.
Environment variable: %q`, authOIDCRolesClaimsEnv)
	f.String(authOIDCRolesClaimsFlag, "", desc)
	_ = viper.BindPFlag(authOIDCRolesClaimsViperKey, f.Lookup(authOIDCRolesClaimsFlag))
}

func NewOIDCConfig(vipr *viper.Viper) *jose.Config {
	issuerURL := vipr.GetString(authOIDCIssuerURLViperKey)

	if issuerURL != "" {
		return jose.NewConfig(
			issuerURL,
			vipr.GetStringSlice(AuthOIDCAudienceViperKey),
			vipr.GetString(authOIDCQKMClaimsViperKey),
			vipr.GetString(authOIDCPermissionsClaimsViperKey),
			vipr.GetString(authOIDCRolesClaimsViperKey),
			5*time.Minute, // TODO: Make the cache ttl an ENV var if needed
		)
	}

	return nil
}
