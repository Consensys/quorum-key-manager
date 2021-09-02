package tls

import (
	"crypto/x509"
)

func VerifyCertificateAuthority(certs []*x509.Certificate, serverName string, rootCAs *x509.CertPool, skipVerify bool) error {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         rootCAs,
	}

	if !skipVerify {
		opts.DNSName = serverName
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
