package certificate_test

import (
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/tls/certificate"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/tls/testutils"
	"github.com/stretchr/testify/assert"
)

func TestX509KeyPair(t *testing.T) {
	tests := []struct {
		desc         string
		certPemBlock string
		privPemBlock string
		expectedErr  bool
	}{
		{
			"RSA",
			testutils.RSACertPEM,
			testutils.RSAKeyPEM,
			false,
		},
		{
			"RSA 1 line",
			testutils.OneLineRSACertPEMA,
			testutils.OneLineRSAKeyPEMA,
			false,
		},
		{
			"RSA cert only",
			testutils.RSACertPEM,
			"",
			false,
		},
		{
			"RSA Key only",
			"",
			testutils.RSAKeyPEM,
			false,
		},
		{
			"ECDSA",
			testutils.ECDSACertPEM,
			testutils.ECDSAKeyPEM,
			false,
		},
		{
			"Mix Cert ECDSA / Key RSA",
			testutils.ECDSACertPEM,
			testutils.RSAKeyPEM,
			true,
		},
		{
			"Mix Cert RSA / Key ECDSA",
			testutils.RSACertPEM,
			testutils.ECDSAKeyPEM,
			true,
		},
		{
			"RSA unmatching",
			testutils.OneLineRSACertPEMB,
			testutils.OneLineRSAKeyPEMA,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := certificate.X509KeyPair([]byte(tt.certPemBlock), []byte(tt.privPemBlock))
			if (err != nil) != tt.expectedErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		desc          string
		pemBlock      string
		typ           string
		expectedErr   bool
		expectedCount int
	}{
		{
			"RSA Certificate PEM block",
			testutils.RSACertPEM,
			"CERTIFICATE",
			false,
			1,
		},
		{
			"RSA Key PEM block",
			testutils.RSAKeyPEM,
			"PRIVATE KEY",
			false,
			1,
		},
		{
			"RSA Cert PEM with missing headers",
			testutils.NoHeaderRSACertPEM,
			"CERTIFICATE",
			false,
			1,
		},
		{
			"valid 1 line block",
			testutils.OneLineRSACertPEMA,
			"CERTIFICATE",
			false,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			certs, err := certificate.Decode([]byte(tt.pemBlock), tt.typ)
			if (err != nil) != tt.expectedErr {
				t.Errorf("ParsePEM() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			assert.Len(t, certs, tt.expectedCount, "Count of of certificates should be corrects")
		})
	}
}
