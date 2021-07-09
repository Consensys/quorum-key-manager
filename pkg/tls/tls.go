package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
)

var (
	// ClientAuthTypes Map of allowed TLS ClientAuthType
	ClientAuthTypes = map[string]tls.ClientAuthType{
		"NoClientCert":               tls.NoClientCert,
		"RequestClientCert":          tls.RequestClientCert,
		"RequireAnyClientCert":       tls.RequireAnyClientCert,
		"VerifyClientCertIfGiven":    tls.VerifyClientCertIfGiven,
		"RequireAndVerifyClientCert": tls.RequireAndVerifyClientCert,
	}

	// Versions map of allowed TLS versions
	Versions = map[string]uint16{
		`VersionTLS10`: tls.VersionTLS10,
		`VersionTLS11`: tls.VersionTLS11,
		`VersionTLS12`: tls.VersionTLS12,
		`VersionTLS13`: tls.VersionTLS13,
	}

	// CurveIDs is a Map of TLS elliptic curves from crypto/tls
	// Available CurveIDs defined at https://godoc.org/crypto/tls#CurveID,
	// also allowing rfc names defined at https://tools.ietf.org/html/rfc8446#section-4.2.7
	CurveIDs = map[string]tls.CurveID{
		`secp256r1`: tls.CurveP256,
		`CurveP256`: tls.CurveP256,
		`secp384r1`: tls.CurveP384,
		`CurveP384`: tls.CurveP384,
		`secp521r1`: tls.CurveP521,
		`CurveP521`: tls.CurveP521,
		`x25519`:    tls.X25519,
		`X25519`:    tls.X25519,
	}

	// CipherSuites Map of TLS CipherSuites from crypto/tls
	// Available CipherSuites defined at https://golang.org/pkg/crypto/tls/#pkg-constants
	CipherSuites = map[string]uint16{
		`TLS_RSA_WITH_RC4_128_SHA`:                      tls.TLS_RSA_WITH_RC4_128_SHA,
		`TLS_RSA_WITH_3DES_EDE_CBC_SHA`:                 tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		`TLS_RSA_WITH_AES_128_CBC_SHA`:                  tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		`TLS_RSA_WITH_AES_256_CBC_SHA`:                  tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		`TLS_RSA_WITH_AES_128_CBC_SHA256`:               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		`TLS_RSA_WITH_AES_128_GCM_SHA256`:               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		`TLS_RSA_WITH_AES_256_GCM_SHA384`:               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		`TLS_ECDHE_ECDSA_WITH_RC4_128_SHA`:              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		`TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA`:          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		`TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA`:          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		`TLS_ECDHE_RSA_WITH_RC4_128_SHA`:                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		`TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA`:           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		`TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`:            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		`TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`:            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		`TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256`:       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		`TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256`:         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		`TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`:         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		`TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`:       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		`TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`:         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		`TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`:       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		`TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305`:          tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		`TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256`:   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		`TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305`:        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		`TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256`: tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		`TLS_AES_128_GCM_SHA256`:                        tls.TLS_AES_128_GCM_SHA256,
		`TLS_AES_256_GCM_SHA384`:                        tls.TLS_AES_256_GCM_SHA384,
		`TLS_CHACHA20_POLY1305_SHA256`:                  tls.TLS_CHACHA20_POLY1305_SHA256,
		`TLS_FALLBACK_SCSV`:                             tls.TLS_FALLBACK_SCSV,
	}
)

// Option configures TLS for an entry point
type Option struct {
	Certificates []*certificate.KeyPair `json:"certificates,omitempty" toml:"certificates,omitempty" yaml:"certificates,omitempty" export:"true"`
	CAs          [][]byte               `json:"clientCAs,omitempty" toml:"clientCAs,omitempty" yaml:"clientCAs,omitempty"`

	NextProtos []string `json:"nextProtos,omitempty" toml:"nextProtos,omitempty" yaml:"nextProtos,omitempty" export:"true"`

	CipherSuites     []string `json:"cipherSuites,omitempty" toml:"cipherSuites,omitempty" yaml:"cipherSuites,omitempty"`
	CurvePreferences []string `json:"curvePreferences,omitempty" toml:"curvePreferences,omitempty" yaml:"curvePreferences,omitempty"`

	ClientAuth string `json:"clientAuthType,omitempty" toml:"clientAuthType,omitempty" yaml:"clientAuthType,omitempty"`

	MinVersion string `json:"minVersion,omitempty" toml:"minVersion,omitempty" yaml:"minVersion,omitempty" export:"true"`
	MaxVersion string `json:"maxVersion,omitempty" toml:"maxVersion,omitempty" yaml:"maxVersion,omitempty" export:"true"`

	ServerName string `json:"serverName,omitempty" toml:"serverName,omitempty" yaml:"serverName,omitempty" export:"true"`

	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty" toml:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`

	PreferServerCipherSuites bool `json:"preferServerCipherSuites,omitempty" toml:"preferServerCipherSuites,omitempty" yaml:"preferServerCipherSuites,omitempty" export:"true"`

	SniStrict bool `json:"sniStrict,omitempty" toml:"sniStrict,omitempty" yaml:"sniStrict,omitempty" export:"true"`
}

func (opt *Option) TLSClientAuth() (tls.ClientAuthType, error) {
	v, ok := ClientAuthTypes[opt.ClientAuth]
	if !ok {
		return tls.NoClientCert, fmt.Errorf("invalid TLS client auth type %q", opt.ClientAuth)
	}
	return v, nil
}

func (opt *Option) TLSMinVersion() (uint16, bool) {
	return tlsVersion(opt.MinVersion)
}

func (opt *Option) TLSMaxVersion() (uint16, bool) {
	return tlsVersion(opt.MaxVersion)
}

func tlsVersion(v string) (uint16, bool) {
	rv, ok := Versions[v]
	return rv, ok
}

func NewConfig(opt *Option) (*tls.Config, error) {
	cfg := &tls.Config{
		ServerName:               opt.ServerName,
		NextProtos:               opt.NextProtos,
		PreferServerCipherSuites: opt.PreferServerCipherSuites,
		InsecureSkipVerify:       opt.InsecureSkipVerify,
	}

	// Load Certificates
	for _, certOpt := range opt.Certificates {
		cert, err := certificate.X509(certOpt)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	// Load Client CAs
	if len(opt.CAs) > 0 {
		pool := x509.NewCertPool()
		for _, ca := range opt.CAs {
			cert, err := certificate.X509KeyPair(ca, nil)
			if err != nil {
				return nil, err
			}

			for _, asn1 := range cert.Certificate {
				c, err := x509.ParseCertificate(asn1)
				if err != nil {
					return nil, err
				}
				pool.AddCert(c)
			}
		}
		cfg.ClientCAs = pool
		cfg.RootCAs = pool
	}

	// Set client Auth type
	if len(opt.ClientAuth) > 0 {
		clientAuth, err := opt.TLSClientAuth()
		if err != nil {
			return nil, err
		}
		cfg.ClientAuth = clientAuth

		// Make sure ClientCAs has been properly set
		if cfg.ClientCAs == nil && (clientAuth == tls.VerifyClientCertIfGiven || clientAuth == tls.RequireAndVerifyClientCert) {
			return nil, fmt.Errorf("invalid ClientAuth: %q, CAFiles is required", opt.ClientAuth)
		}
	}

	var ok bool
	if cfg.MinVersion, ok = opt.TLSMinVersion(); ok {
		cfg.PreferServerCipherSuites = true
	}
	if cfg.MaxVersion, ok = opt.TLSMaxVersion(); ok {
		cfg.PreferServerCipherSuites = true
	}

	// Set the list of CipherSuites if set in the config
	if len(opt.CipherSuites) > 0 {
		for _, cipher := range opt.CipherSuites {
			if suite, ok := CipherSuites[cipher]; ok {
				cfg.CipherSuites = append(cfg.CipherSuites, suite)
			} else {
				// CipherSuite listed in the toml does not exist in our listed
				return nil, fmt.Errorf("invalid CipherSuite: %s", cipher)
			}
		}
	}

	// Set the list of CurvePreferences/CurveIDs if set in the config
	if len(opt.CurvePreferences) > 0 {
		// if our list of CurvePreferences/CurveIDs is defined in the config, we can re-initialize the list as empty
		for _, curve := range opt.CurvePreferences {
			if curveID, ok := CurveIDs[curve]; ok {
				cfg.CurvePreferences = append(cfg.CurvePreferences, curveID)
			} else {
				// CurveID listed in the toml does not exist in our listed
				return nil, fmt.Errorf("invalid CurvePreference: %s", curve)
			}
		}
	}

	return cfg, nil
}
