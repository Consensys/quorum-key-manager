package flags

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	authtls "github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"

	"github.com/consensys/quorum-key-manager/pkg/jwt"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/auth"
	apikey "github.com/consensys/quorum-key-manager/src/auth/authenticator/api-key"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authOIDCCACertFileViperKey, authOIDCCACertFileEnv)
	_ = viper.BindEnv(AuthOIDCCAKeyFileViperKey, authOIDCCAKeyFileEnv)
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)

	viper.SetDefault(authOIDCClaimUsernameViperKey, authOIDCClaimUsernameDefault)
	_ = viper.BindEnv(authOIDCClaimUsernameViperKey, authOIDCClaimUsernameEnv)
	viper.SetDefault(authOIDCClaimGroupViperKey, authOIDCClaimGroupDefault)
	_ = viper.BindEnv(authOIDCClaimGroupViperKey, authOIDCClaimGroupEnv)

	viper.SetDefault(authAPIKeyFileViperKey, authAPIKeyDefaultFileFlag)
	_ = viper.BindEnv(authAPIKeyFileViperKey, authAPIKeyFileEnv)

	_ = viper.BindEnv(authTLSCertsFileViperKey, authTLSCertsFileEnv)

}

const (
	csvSeparator      = ';'
	csvCommentsMarker = '#'
	csvRowLen         = 3
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
	authOIDCCACertFileFlag     = "auth-oidc-ca-cert"
	authOIDCCACertFileViperKey = "auth.oidc.ca.cert"
	authOIDCCACertFileDefault  = ""
	authOIDCCACertFileEnv      = "AUTH_OIDC_CA_CERT"
)

const (
	authOIDCIssuerURLFlag     = "auth-oidc-issuer-url"
	authOIDCIssuerURLViperKey = "auth.oidc.issuer.url"
	authOIDCIssuerURLDefault  = ""
	authOIDCIssuerURLEnv      = "AUTH_OIDC_ISSUER_URL"
)

const (
	authOIDCCAKeyFileFlag     = "auth-oidc-ca-key"
	AuthOIDCCAKeyFileViperKey = "auth.oidc.ca.key"
	authOIDCCAKeyFileDefault  = ""
	authOIDCCAKeyFileEnv      = "AUTH_OIDC_CA_KEY"
)

const (
	authOIDCClaimUsernameFlag     = "auth-oidc-claim-username"
	authOIDCClaimUsernameViperKey = "auth.oidc.claim.username"
	authOIDCClaimUsernameDefault  = "sub"
	authOIDCClaimUsernameEnv      = "AUTH_OIDC_CLAIM_USERNAME"
)

const (
	authOIDCClaimGroupFlag     = "auth-oidc-claim-groups"
	authOIDCClaimGroupViperKey = "auth.oidc.claim.groups"
	authOIDCClaimGroupDefault  = "scope"
	authOIDCClaimGroupEnv      = "AUTH_OIDC_CLAIM_GROUPS"
)

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

func AuthFlags(f *pflag.FlagSet) {
	authOIDCCAFile(f)
	authOIDCIssuerServer(f)
	AuthOIDCClaimUsername(f)
	AuthOIDCClaimGroups(f)
	authTLSCertFile(f)
	authAPIKeyFile(f)
}

// Use only on generate-token utils
func AuthOIDCCertKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect CA Cert filepath.
Environment variable: %q`, authOIDCCAKeyFileEnv)
	f.String(authOIDCCAKeyFileFlag, authOIDCCAKeyFileDefault, desc)
	_ = viper.BindPFlag(AuthOIDCCAKeyFileViperKey, f.Lookup(authOIDCCAKeyFileFlag))
}

func authOIDCCAFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect CA Key filepath.
Environment variable: %q`, authOIDCCACertFileEnv)
	f.String(authOIDCCACertFileFlag, authOIDCCACertFileDefault, desc)
	_ = viper.BindPFlag(authOIDCCACertFileViperKey, f.Lookup(authOIDCCACertFileFlag))
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
	f.String(authOIDCClaimUsernameFlag, authOIDCClaimUsernameDefault, desc)
	_ = viper.BindPFlag(authOIDCClaimUsernameViperKey, f.Lookup(authOIDCClaimUsernameFlag))
}

func AuthOIDCClaimGroups(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for groups.
Environment variable: %q`, authOIDCClaimGroupEnv)
	f.String(authOIDCClaimGroupFlag, authOIDCClaimGroupDefault, desc)
	_ = viper.BindPFlag(authOIDCClaimGroupViperKey, f.Lookup(authOIDCClaimGroupFlag))
}

func NewAuthConfig(vipr *viper.Viper) (*auth.Config, error) {
	// OIDC
	certsOIDC := []*x509.Certificate{}

	fileCertOIDC, err := oidcCert(vipr)
	if err != nil {
		return nil, err
	} else if fileCertOIDC != nil {
		certsOIDC = append(certsOIDC, fileCertOIDC)
	}

	issuerCerts, err := oidcIssuerURL(vipr)
	if err != nil {
		return nil, err
	} else if issuerCerts != nil {
		certsOIDC = append(certsOIDC, issuerCerts...)
	}

	oidcCfg := oidc.NewConfig(vipr.GetString(authOIDCClaimUsernameViperKey),
		vipr.GetString(authOIDCClaimGroupViperKey), certsOIDC...)

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
		OIDC:   oidcCfg,
		APIKEY: apiKeyCfg,
		TLS:    tlsCfg,
	}, nil

}

func oidcCert(vipr *viper.Viper) (*x509.Certificate, error) {
	caFile := vipr.GetString(authOIDCCACertFileViperKey)
	if caFile == "" {
		return nil, nil
	}

	_, err := os.Stat(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file. %s", err.Error())
	}

	caFileContent, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	bCert, err := certificate.Decode(caFileContent, "CERTIFICATE")
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(bCert[0])
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func oidcIssuerURL(vipr *viper.Viper) ([]*x509.Certificate, error) {
	issuerServer := vipr.GetString(authOIDCIssuerURLViperKey)
	if issuerServer == "" {
		return nil, nil
	}

	jwks, err := jwt.RetrieveKeySet(context.Background(), http.DefaultClient, issuerServer)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth server jwks: %s", issuerServer)
	}

	certs := []*x509.Certificate{}
	for _, kw := range jwks.Keys {
		certs = append(certs, kw.Certificates...)
	}

	return certs, nil
}

func apiKeyCsvFile(vipr *viper.Viper) (map[string]apikey.UserClaims, error) {
	// Open the file
	csvFileName := vipr.GetString(authAPIKeyFileViperKey)
	if csvFileName == "" {
		return nil, nil
	}
	csvfile, err := os.Open(csvFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read api-key filepath %s: %w", csvFileName, err)
	}

	defer func(csvfile *os.File) {
		_ = csvfile.Close()
	}(csvfile)

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

		retFile[cells[0]] = apikey.UserClaims{UserName: cells[1],
			Claims: strings.Split(cells[2], ","),
		}
	}

	return retFile, nil
}

func tlsAuthCerts(vipr *viper.Viper) ([]*x509.Certificate, error) {
	caFile := vipr.GetString(authTLSCertsFileViperKey)
	if caFile == "" {
		return nil, nil
	}
	_, err := os.Stat(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file. %s", err.Error())
	}

	caFileContent, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	bCert, err := certificate.Decode(caFileContent, "CERTIFICATE")
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(bCert[0])
	if err != nil {
		return nil, err
	}

	return []*x509.Certificate{cert}, nil
}
