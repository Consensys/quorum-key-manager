package flags

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	auth2 "github.com/consensys/quorum-key-manager/pkg/auth"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oicd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authOICDCACertFileViperKey, authOICDCACertFileEnv)
	_ = viper.BindEnv(AuthOICDCAKeyFileViperKey, authOICDCAKeyFileEnv)
	_ = viper.BindEnv(authOICDIssuerURLViperKey, authOICDIssuerURLEnv)

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
	authOICDIssuerURLFlag     = "auth-oidc-issuer-url"
	authOICDIssuerURLViperKey = "auth.oidc.issuer.url"
	authOICDIssuerURLDefault  = ""
	authOICDIssuerURLEnv      = "AUTH_OICD_ISSUER_URL"
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
	authOICDIssuerServer(f)
	AuthOICDClaimUsername(f)
	AuthOICDClaimGroups(f)
}

// Use only on generate-token utils 
func AuthOICDCertKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect CA Cert filepath.
Environment variable: %q`, authOICDCAKeyFileEnv)
	f.String(authOICDCAKeyFileFlag, authOICDCAKeyFileDefault, desc)
	_ = viper.BindPFlag(AuthOICDCAKeyFileViperKey, f.Lookup(authOICDCAKeyFileFlag))
}

func authOICDCAFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect CA Key filepath.
Environment variable: %q`, authOICDClaimUsernameEnv)
	f.String(authOICDClaimUsernameFlag, authOICDClaimUsernameDefault, desc)
	_ = viper.BindPFlag(authOICDClaimUsernameViperKey, f.Lookup(authOICDClaimUsernameFlag))
}

func authOICDIssuerServer(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect issuer server domain (ie. https://quorum-key-manager.eu.auth0.com/.well-known/jwks.json).
Environment variable: %q`, authOICDIssuerURLEnv)
	f.String(authOICDIssuerURLFlag, authOICDIssuerURLDefault, desc)
	_ = viper.BindPFlag(authOICDIssuerURLViperKey, f.Lookup(authOICDIssuerURLFlag))
}

func AuthOICDClaimUsername(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for username.
Environment variable: %q`, authOICDClaimGroupEnv)
	f.String(authOICDClaimGroupFlag, authOICDClaimGroupDefault, desc)
	_ = viper.BindPFlag(authOICDClaimGroupViperKey, f.Lookup(authOICDClaimGroupFlag))
}

func AuthOICDClaimGroups(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Token path claims for groups.
Environment variable: %q`, authOICDCACertFileEnv)
	f.String(authOICDCACertFileFlag, authOICDCACertFileDefault, desc)
	_ = viper.BindPFlag(authOICDCACertFileViperKey, f.Lookup(authOICDCACertFileFlag))
}

func NewAuthConfig(vipr *viper.Viper) (*auth.Config, error) {
	var oicdCfg = &oicd.Config{}
	certs := []*x509.Certificate{}
	
	fileCert, err := fileCertificate(vipr)
	if err != nil {
		return nil, err
	} else if fileCert != nil {
		certs = append(certs, fileCert)
	}
	
	issuerCerts, err := issuerCertificates(vipr)
	if err != nil {
		return nil, err
	} else if issuerCerts != nil {
		certs = append(certs, issuerCerts...)
	}

	oicdCfg = oicd.NewConfig(vipr.GetString(authOICDClaimUsernameViperKey), 
			vipr.GetString(authOICDClaimGroupViperKey), certs...)
	
	return &auth.Config{OICD: oicdCfg}, nil
}

func fileCertificate(vipr *viper.Viper) (*x509.Certificate, error) {
	caFile := vipr.GetString(authOICDCACertFileViperKey)
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

func issuerCertificates(vipr *viper.Viper) ([]*x509.Certificate, error) {
	issuerServer := vipr.GetString(authOICDIssuerURLViperKey)
	if issuerServer == "" {
		return nil, nil
	}
	
	jwks, err := auth2.RetrieveKeySet(context.Background(), http.DefaultClient, issuerServer)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth server jwks", issuerServer)
	}

	certs := []*x509.Certificate{}
	for _, kw := range jwks.Keys {
		certs = append(certs, kw.Certificates...)
	}
	
	return certs, nil
}

