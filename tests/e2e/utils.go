package e2e

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
)

const MaxRetries = 5

type callFunc func() error
type logFunc func(format string, args ...interface{})

func retryOn(call callFunc, logger logFunc, errMsg string, httpStatusCode, retries int) error {
	for {
		err := call()
		if httpError, ok := err.(*client.ResponseError); retries <= 0 || !ok || httpError.StatusCode != httpStatusCode {
			if err != nil {
				return err
			}
			break
		}

		logger("%s (retrying in 1 second...)", errMsg)
		time.Sleep(time.Second)
		retries--
	}

	return nil
}

func generateAPIKey(key string) string {
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func generateClientCert(certFile, keyFile string) (*tls.Certificate, error) {
	curDir, _ := os.Getwd()
	cert, err := tls.LoadX509KeyPair(path.Join(curDir, certFile), path.Join(curDir, keyFile))
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

type testHttpTransport struct {
	token  string
	apiKey string
	cert   *tls.Certificate
}

func NewTestHttpTransport(token, apiKey string, cert *tls.Certificate) http.RoundTripper {
	return &testHttpTransport{
		token:  token,
		apiKey: apiKey,
		cert:   cert,
	}
}

func (t *testHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	defaultTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	switch {
	case t.cert != nil:
		defaultTransport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{*t.cert},
				GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return t.cert, nil
				},
			},
		}
	case t.token != "":
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	case t.apiKey != "":
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", t.apiKey))
	}

	return defaultTransport.RoundTrip(req)
}
