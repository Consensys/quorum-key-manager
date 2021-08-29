package e2e

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/consensys/quorum-key-manager/cmd/flags"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/jwt"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/spf13/viper"
)

const MAX_RETRIES = 5

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

func generateJWT(keyFile, scope, sub string) (string, error) {
	curDir, _ := os.Getwd()
	keyFileContent, err := ioutil.ReadFile(path.Join(curDir, keyFile))
	if err != nil {
		return "", err
	}

	var keys [][]byte
	keys, err = certificate.Decode(keyFileContent, "PRIVATE KEY")
	if err != nil {
		return "", err
	}
	certKey, err := certificate.ParsePrivateKey(keys[0])
	if err != nil {
		return "", err
	}
	generator, err := jwt.NewTokenGenerator(certKey)
	if err != nil {
		return "", err
	}

	cfgSubject := viper.GetString(flags.AuthOIDCClaimUsernameViperKey)
	cfgScope := viper.GetString(flags.AuthOIDCClaimGroupViperKey)

	return generator.GenerateAccessToken(map[string]interface{}{
		cfgSubject: sub,
		cfgScope:   scope,
	}, time.Hour)
}

type TestHttpTransport struct {
	token            string
	defaultTransport http.RoundTripper
}

func NewTestHttpTransport(token string) http.RoundTripper {
	return &TestHttpTransport{
		token: token,
		defaultTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (t *TestHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}

	return t.defaultTransport.RoundTrip(req)
}
