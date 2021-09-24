package types

import "time"

// TODO: Keep this type as an example of the request type when a manifest is read
type CreateHashicorpStoreRequest struct {
	MountPoint    string        `json:"mountPoint" validate:"required" example:"secret"`
	Address       string        `json:"address"  validate:"required" example:"https://hashicorp:8200"`
	Token         string        `json:"token,omitempty" example:"s.W7IMlFuBGsTaR6uHLcGDw9Mq"`
	TokenPath     string        `json:"tokenPath,omitempty" example:"/vault/token/.my_token"`
	Namespace     string        `json:"namespace,omitempty" example:"namespace"`
	CACert        string        `json:"CACert,omitempty" example:"/vault/tls/ca.crt"`
	CAPath        string        `json:"CAPath,omitempty" example:"/vault/tls/root-certs"`
	ClientCert    string        `json:"clientCert,omitempty" example:"/vault/tls/client.crt"`
	ClientKey     string        `json:"clientKey,omitempty" example:"/vault/tls/client.key"`
	TLSServerName string        `json:"TLSServerName,omitempty" example:"server-name"`
	ClientTimeout time.Duration `json:"clientTimeout,omitempty" example:"60s"`
	RateLimit     float64       `json:"rateLimit,omitempty" example:"0"`
	BurstLimit    int           `json:"burstLimit,omitempty" example:"0"`
	MaxRetries    int           `json:"maxRetries,omitempty" example:"2"`
	SkipVerify    bool          `json:"skipVerify,omitempty" example:"false"`
}
