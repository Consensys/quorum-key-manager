package hashicorp

import (
	"crypto/tls"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	hashicorp "github.com/hashicorp/vault/api"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type Config struct {
	TokenFilePath string
	MountPoint    string
	SecretPath    string
	RateLimit     float64
	BurstLimit    int
	Address       string
	CACert        string
	CAPath        string
	ClientCert    string
	ClientKey     string
	ClientTimeout time.Duration
	MaxRetries    int
	SkipVerify    bool
	TLSServerName string
	Namespace     string
}

func (cfg *Config) ToVaultConfig() *hashicorp.Config {
	// Create Vault Configuration
	config := hashicorp.DefaultConfig()
	config.Address = cfg.Address
	config.HttpClient = cleanhttp.DefaultClient()
	config.HttpClient.Timeout = time.Second * 60

	// Create Transport
	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		config.Error = err
		return config
	}

	// Configure TLS
	tlsConfig := &hashicorp.TLSConfig{
		CACert:        cfg.CACert,
		CAPath:        cfg.CAPath,
		ClientCert:    cfg.ClientCert,
		ClientKey:     cfg.ClientKey,
		TLSServerName: cfg.TLSServerName,
		Insecure:      cfg.SkipVerify,
	}

	_ = config.ConfigureTLS(tlsConfig)

	config.Limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.BurstLimit)
	config.MaxRetries = cfg.MaxRetries
	config.Timeout = cfg.ClientTimeout

	// Ensure redirects are not automatically followed
	// Note that this is sane for the API client as it has its own
	// redirect handling logic (and thus also for command/meta),
	// but in e.g. http_test actual redirect handling is necessary
	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Returning this value causes the Go net library to not close the
		// response body and to nil out the error. Otherwise retry clients may
		// try three times on every redirect because it sees an error from this
		// function (to prevent redirects) passing through to it.
		return http.ErrUseLastResponse
	}

	config.Backoff = retryablehttp.LinearJitterBackoff

	return config
}
