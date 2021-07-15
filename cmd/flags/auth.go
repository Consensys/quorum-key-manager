package flags

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/consensys/quorum-key-manager/pkg/jwt"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	authTLSCertsFileFlag        = "auth-tls-client-certs"
	authTLSCertsFileViperKey    = "auth.tls.client.certs"
	authTLSCertsDefaultFileFlag = ""
	authTLSCertsCertsFileEnv    = "AUTH_TLS_CLIENT_CERTS"
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
	authOIDCClaimUsernameDefault  = "qkm.auth.username"
	authOIDCClaimUsernameEnv      = "AUTH_OIDC_CLAIM_USERNAME"
)

const (
	authOIDCClaimGroupFlag     = "auth-oidc-claim-groups"
	authOIDCClaimGroupViperKey = "auth.oidc.claim.groups"
	authOIDCClaimGroupDefault  = "qkm.auth.groups"
	authOIDCClaimGroupEnv      = "AUTH_OIDC_CLAIM_GROUPS"
)

func init() {
	_ = viper.BindEnv(authOIDCCACertFileViperKey, authOIDCCACertFileEnv)
	_ = viper.BindEnv(AuthOIDCCAKeyFileViperKey, authOIDCCAKeyFileEnv)
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)

	viper.SetDefault(authOIDCClaimUsernameViperKey, authOIDCClaimUsernameDefault)
	_ = viper.BindEnv(authOIDCClaimUsernameViperKey, authOIDCClaimUsernameEnv)

	viper.SetDefault(authOIDCClaimGroupViperKey, authOIDCClaimGroupDefault)
	_ = viper.BindEnv(authOIDCClaimGroupViperKey, authOIDCClaimGroupEnv)

	_ = viper.BindEnv(authTLSCertsFileViperKey, authTLSCertsCertsFileEnv)

}

// Use only on generate-token utils
func AuthTLSCertKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect CA Cert filepath.
Environment variable: %q`, authTLSCertsCertsFileEnv)
	f.String(authTLSCertsFileFlag, authTLSCertsDefaultFileFlag, desc)
	_ = viper.BindPFlag(authTLSCertsFileViperKey, f.Lookup(authTLSCertsFileFlag))
}

func clientCertificate(vipr *viper.Viper) (*x509.Certificate, error) {
	caFile := vipr.GetString(authTLSCertsFileViperKey)
	_, err := os.Stat(caFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read CA file. %s", err.Error())
		}
		return nil, nil
	}

	caFileContent, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	bCert, err := certificate.Decode(caFileContent, "CERTIFICATE")
	cert, err := x509.ParseCertificate(bCert[0])
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func AuthFlags(f *pflag.FlagSet) {
	authOIDCCAFile(f)
	authOIDCIssuerServer(f)
	AuthOIDCClaimUsername(f)
	AuthOIDCClaimGroups(f)
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
Environment variable: %q`, authOIDCClaimUsernameEnv)
	f.String(authOIDCClaimUsernameFlag, authOIDCClaimUsernameDefault, desc)
	_ = viper.BindPFlag(authOIDCClaimUsernameViperKey, f.Lookup(authOIDCClaimUsernameFlag))
}

func authOIDCIssuerServer(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect issuer server domain (ie. https://quorum-key-manager.eu.auth0.com/.well-known/jwks.json).
Environment variable: %q`, authOIDCIssuerURLEnv)
	f.String(authOIDCIssuerURLFlag, authOIDCIssuerURLDefault, desc)
	_ = viper.BindPFlag(authOIDCIssuerURLViperKey, f.Lookup(authOIDCIssuerURLFlag))
}

func AuthOIDCClaimUsername(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for username.
Environment variable: %q`, authOIDCClaimGroupEnv)
	f.String(authOIDCClaimGroupFlag, authOIDCClaimGroupDefault, desc)
	_ = viper.BindPFlag(authOIDCClaimGroupViperKey, f.Lookup(authOIDCClaimGroupFlag))
}

func AuthOIDCClaimGroups(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for groups.
Environment variable: %q`, authOIDCCACertFileEnv)
	f.String(authOIDCCACertFileFlag, authOIDCCACertFileDefault, desc)
	_ = viper.BindPFlag(authOIDCCACertFileViperKey, f.Lookup(authOIDCCACertFileFlag))
}

func NewAuthConfig(vipr *viper.Viper) (*auth.Config, error) {
	// OIDC part
	certsOIDC := []*x509.Certificate{}

	fileCertOIDC, err := fileCertificate(vipr)
	if err != nil {
		return nil, err
	} else if fileCertOIDC != nil {
		certsOIDC = append(certsOIDC, fileCertOIDC)
	}

	issuerCerts, err := issuerCertificates(vipr)
	if err != nil {
		return nil, err
	} else if issuerCerts != nil {
		certsOIDC = append(certsOIDC, issuerCerts...)
	}

	oidcCfg := oidc.NewConfig(vipr.GetString(authOIDCClaimUsernameViperKey),
		vipr.GetString(authOIDCClaimGroupViperKey), certsOIDC...)

	// TLS part
	var tlsCfg = &tls.Config{}
	certsTLS := []*x509.Certificate{}

	fileCertTLS, err := clientCertificate(vipr)
	if err != nil {
		return nil, err
	} else if fileCertTLS != nil {
		certsTLS = append(certsTLS, fileCertTLS)
	}

	tlsCfg = tls.NewConfig(certsTLS...)

	return &auth.Config{OIDC: oidcCfg, TLS: tlsCfg}, nil

}

func fileCertificate(vipr *viper.Viper) (*x509.Certificate, error) {
	caFile := vipr.GetString(authOIDCCACertFileViperKey)
	_, err := os.Stat(caFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read CA file. %s", err.Error())
		}
		return nil, nil
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

func issuerCertificates(vipr *viper.Viper) ([]*x509.Certificate, error) {
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
