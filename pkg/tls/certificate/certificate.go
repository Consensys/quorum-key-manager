package certificate

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

var (
	errKeyPairTypes = fmt.Errorf("private key type does not match public key type")
	errKeyPair      = fmt.Errorf("private key does not match public key")
)

type KeyPair struct {
	Cert []byte `json:"cert,omitempty" toml:"cert,omitempty" yaml:"cert,omitempty"`
	Key  []byte `json:"key,omitempty" toml:"key,omitempty" yaml:"key,omitempty"`
}

func X509(pair *KeyPair) (tls.Certificate, error) {
	return X509KeyPair(pair.Cert, pair.Key)
}

func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (tls.Certificate, error) {
	var (
		cert tls.Certificate
		err  error
	)
	// Parse CERTIFICATE if provided
	if len(certPEMBlock) > 0 {
		cert.Certificate, err = Decode(certPEMBlock, "CERTIFICATE")
		if err != nil {
			return tls.Certificate{}, err
		}

		// Set certificate leaf
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return tls.Certificate{}, err
		}
	}

	// Parse PRIVATE KEY if provided
	if len(keyPEMBlock) > 0 {
		var keys [][]byte
		keys, err = Decode(keyPEMBlock, "PRIVATE KEY")
		if err != nil {
			return tls.Certificate{}, err
		}

		cert.PrivateKey, err = parsePrivateKey(keys[0])
		if err != nil {
			return tls.Certificate{}, err
		}
	}

	// Test certificate and private key matches
	if cert.Leaf != nil && cert.PrivateKey != nil {
		switch pub := cert.Leaf.PublicKey.(type) {
		case *rsa.PublicKey:
			priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
			if !ok {
				return tls.Certificate{}, errKeyPairTypes
			}
			if pub.N.Cmp(priv.N) != 0 {
				return tls.Certificate{}, errKeyPair
			}
		case *ecdsa.PublicKey:
			priv, ok := cert.PrivateKey.(*ecdsa.PrivateKey)
			if !ok {
				return tls.Certificate{}, errKeyPairTypes
			}
			if pub.X.Cmp(priv.X) != 0 || pub.Y.Cmp(priv.Y) != 0 {
				return tls.Certificate{}, errKeyPair
			}
		case ed25519.PublicKey:
			priv, ok := cert.PrivateKey.(ed25519.PrivateKey)
			if !ok {
				return tls.Certificate{}, errKeyPairTypes
			}
			if !bytes.Equal(priv.Public().(ed25519.PublicKey), pub) {
				return tls.Certificate{}, errKeyPair
			}
		default:
			return tls.Certificate{}, fmt.Errorf("unknown public key algorithm")
		}
	}

	return cert, nil
}

func Decode(block []byte, typ string) ([][]byte, error) {
	var (
		certs [][]byte
		err   error
	)
	// Parses assuming block has a valid  format
	if certs, err = decode(block, typ); err == nil {
		return certs, nil
	}
	rErr := err

	// Parses assuming block headers are missing
	if certs, err = decode([]byte(fmt.Sprintf("-----BEGIN %v-----\n%v\n-----END %v-----", typ, string(block), typ)), typ); err == nil {
		return certs, nil
	}

	// Parses assuming block is 1 line
	if block, err = decodeBase64(block); err != nil {
		return nil, rErr
	}

	if certs, err := decode([]byte(fmt.Sprintf("-----BEGIN %v-----\n%v\n-----END %v-----", typ, string(block), typ)), typ); err == nil {
		return certs, nil
	}

	return nil, rErr
}

func decode(raw []byte, typ string) ([][]byte, error) {
	var (
		blocks            [][]byte
		skippedBlockTypes []string
	)
	for {
		var block *pem.Block
		block, raw = pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == typ || strings.HasSuffix(block.Type, " "+typ) {
			blocks = append(blocks, block.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, block.Type)
		}
	}

	if len(blocks) == 0 {
		if len(skippedBlockTypes) == 0 {
			return nil, fmt.Errorf("failed to find any  data in input")
		}
		return nil, fmt.Errorf("failed to find %q  block in input after skipping  blocks of the following types: %v", typ, skippedBlockTypes)
	}

	return blocks, nil
}

func decodeBase64(data []byte) ([]byte, error) {
	resultData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(resultData, data)
	if err != nil {
		return nil, err
	}

	resultData = resultData[:n]

	return resultData, nil
}

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("found unknown private key type in PKCS#8 wrapping")
		}
	}

	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("failed to parse private key")
}
