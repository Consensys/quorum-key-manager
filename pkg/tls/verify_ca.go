package tls

import (
	"crypto/tls"
	"crypto/x509"
)

func VerifyCertificateAuthority(conn *tls.Conn, tlsConf *tls.Config) error {
	if err := conn.Handshake(); err != nil {
		return err
	}

	certs := conn.ConnectionState().PeerCertificates
	opts := x509.VerifyOptions{
		DNSName:       conn.ConnectionState().ServerName,
		Intermediates: x509.NewCertPool(),
		Roots:         tlsConf.RootCAs,
	}

	for i, cert := range certs {
		if i == 0 {
			continue
		}
		opts.Intermediates.AddCert(cert)
	}

	_, err := certs[0].Verify(opts)

	return err
}
