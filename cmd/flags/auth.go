package flags

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/http/middlewares/jwt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/consensys/quorum-key-manager/src/auth"
	apikey "github.com/consensys/quorum-key-manager/src/auth/authenticator/api-key"
	authtls "github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(AuthOIDCPrivKeyViperKey, authOIDCPrivKeyEnv)
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)

	viper.SetDefault(AuthOIDCClaimUsernameViperKey, AuthOIDCClaimUsernameDefault)
	_ = viper.BindEnv(AuthOIDCClaimUsernameViperKey, authOIDCClaimUsernameEnv)
	viper.SetDefault(AuthOIDCClaimPermissionsViperKey, AuthOIDCClaimPermissionsDefault)
	_ = viper.BindEnv(AuthOIDCClaimPermissionsViperKey, authOIDCClaimPermissionsEnv)
	viper.SetDefault(authOIDCClaimRolesViperKey, AuthOIDCClaimRolesDefault)
	_ = viper.BindEnv(authOIDCClaimRolesViperKey, authOIDCClaimRolesEnv)

	viper.SetDefault(authAPIKeyFileViperKey, authAPIKeyDefaultFileFlag)
	_ = viper.BindEnv(authAPIKeyFileViperKey, authAPIKeyFileEnv)

	_ = viper.BindEnv(authTLSCertsFileViperKey, authTLSCertsFileEnv)
}

const (
	csvSeparator         = ','
	csvCommentsMarker    = '#'
	csvRowLen            = 4
	csvHashOffset        = 0
	csvUserOffset        = 1
	csvPermissionsOffset = 2
	csvRolesOffset       = 3
)

const (
	authAPIKeyFileFlag        = "auth-api-key-file"
	authAPIKeyFileViperKey    = "auth.api.key.file"
	authAPIKeyDefaultFileFlag = ""
	authAPIKeyFileEnv         = "AUTH_API_KEY_FILE"
)

const (
	authTLSCertsFileFlag     = "auth-tls-ca"
	authTLSCertsFileViperKey = "auth.tls.ca"
	authTLSCertsFileDefault  = ""
	authTLSCertsFileEnv      = "AUTH_TLS_CA"
)

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
	authOIDCClaimUsernameFlag     = "auth-oidc-claim-username"
	AuthOIDCClaimUsernameViperKey = "auth.oidc.claim.username"
	AuthOIDCClaimUsernameDefault  = "sub"
	authOIDCClaimUsernameEnv      = "AUTH_OIDC_CLAIM_USERNAME"
)

const (
	authOIDCClaimPermissionsFlag     = "auth-oidc-claim-permissions"
	AuthOIDCClaimPermissionsViperKey = "auth.oidc.claim.permissions"
	AuthOIDCClaimPermissionsDefault  = "scope"
	authOIDCClaimPermissionsEnv      = "AUTH_OIDC_CLAIM_PERMISSIONS"
)

const (
	authOIDCClaimRolesFlag     = "auth-oidc-claim-roles"
	authOIDCClaimRolesViperKey = "auth.oidc.claim.roles"
	AuthOIDCClaimRolesDefault  = "qkm.roles"
	authOIDCClaimRolesEnv      = "AUTH_OIDC_CLAIM_ROLES"
)

func AuthFlags(f *pflag.FlagSet) {
	authOIDCIssuerServer(f)
	AuthOIDCClaimUsername(f)
	AuthOIDCClaimPermissions(f)
	AuthOIDCClaimRoles(f)
	authTLSCertFile(f)
	authAPIKeyFile(f)
}

func authOIDCIssuerServer(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect issuer server domain (ie. https://quorum-key-manager.eu.auth0.com/.well-known/jwks.json).
Environment variable: %q`, authOIDCIssuerURLEnv)
	f.String(authOIDCIssuerURLFlag, authOIDCIssuerURLDefault, desc)
	_ = viper.BindPFlag(authOIDCIssuerURLViperKey, f.Lookup(authOIDCIssuerURLFlag))
}

func AuthOIDCClaimUsername(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for username.
Environment variable: %q`, authOIDCClaimUsernameEnv)
	f.String(authOIDCClaimUsernameFlag, AuthOIDCClaimUsernameDefault, desc)
	_ = viper.BindPFlag(AuthOIDCClaimUsernameViperKey, f.Lookup(authOIDCClaimUsernameFlag))
}

func AuthOIDCClaimPermissions(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for permissions.
Environment variable: %q`, authOIDCClaimPermissionsEnv)
	f.String(authOIDCClaimPermissionsFlag, AuthOIDCClaimPermissionsDefault, desc)
	_ = viper.BindPFlag(AuthOIDCClaimPermissionsViperKey, f.Lookup(authOIDCClaimPermissionsFlag))
}

func AuthOIDCClaimRoles(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for roles.
Environment variable: %q`, authOIDCClaimRolesEnv)
	f.String(authOIDCClaimRolesFlag, AuthOIDCClaimRolesDefault, desc)
	_ = viper.BindPFlag(authOIDCClaimRolesViperKey, f.Lookup(authOIDCClaimRolesFlag))
}

func authTLSCertFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Authenticator Cert filepath.
Environment variable: %q`, authTLSCertsFileEnv)
	f.String(authTLSCertsFileFlag, authTLSCertsFileDefault, desc)
	_ = viper.BindPFlag(authTLSCertsFileViperKey, f.Lookup(authTLSCertsFileFlag))
}

func authAPIKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Authenticator Cert filepath.
Environment variable: %q`, authAPIKeyFileEnv)
	f.String(authAPIKeyFileFlag, authAPIKeyDefaultFileFlag, desc)
	_ = viper.BindPFlag(authAPIKeyFileViperKey, f.Lookup(authAPIKeyFileFlag))
}

func NewAuthConfig(vipr *viper.Viper) (*auth.Config, error) {
	// OIDC
	oidcCfg := jwt.NewConfig(
		vipr.GetString(authOIDCIssuerURLViperKey),
		vipr.GetString(AuthOIDCClaimUsernameViperKey),
		vipr.GetString(AuthOIDCClaimPermissionsViperKey),
		vipr.GetString(authOIDCClaimRolesViperKey),
		time.Minute, // TODO: Make the cache ttl an ENV var if needed
	)

	// API-KEY
	var apiKeyCfg = &apikey.Config{}
	fileAPIKeys, err := apiKeyCsvFile(vipr)
	if err != nil {
		return nil, err
	} else if fileAPIKeys != nil {
		apiKeyCfg = apikey.NewConfig(fileAPIKeys, base64.StdEncoding, sha256.New())
	}

	// TLS
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

func apiKeyCsvFile(vipr *viper.Viper) (map[string]apikey.UserClaims, error) {
	// Open the file
	csvFileName := vipr.GetString(authAPIKeyFileViperKey)
	if csvFileName == "" {
		return nil, nil
	}
	csvfile, err := os.Open(csvFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read api-key csv file '%s': %w", csvFileName, err)
	}

	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	// Set separator
	r.Comma = csvSeparator
	// ignore comments in file
	r.Comment = csvCommentsMarker

	retFile := make(map[string]apikey.UserClaims)

	// Iterate through the lines
	for {
		// Read each line from csv
		cells, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cells) != csvRowLen {
			return nil, fmt.Errorf("invalid number of cells in file %s should be %d", csvfile.Name(), csvRowLen)
		}

		retFile[cells[csvHashOffset]] = apikey.UserClaims{
			UserName:    cells[csvUserOffset],
			Permissions: strings.Split(cells[csvPermissionsOffset], ","),
			Roles:       strings.Split(cells[csvRolesOffset], ","),
		}
	}

	return retFile, nil
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
