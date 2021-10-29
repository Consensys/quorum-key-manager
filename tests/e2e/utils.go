package e2e

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"io/ioutil"
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

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getJWT(idpURL, clientID, clientSecret, audience string) (string, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	})

	resp, err := http.DefaultClient.Post(idpURL, jsonrpc.ContentType, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	acessToken := &accessTokenResponse{}
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(acessToken); err != nil {
			return "", err
		}

		return acessToken.AccessToken, nil
	}

	// Read body
	respMsg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf(string(respMsg))
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
	case t.apiKey != "":
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(t.apiKey))))
	case t.token != "":
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}

	return defaultTransport.RoundTrip(req)
}
