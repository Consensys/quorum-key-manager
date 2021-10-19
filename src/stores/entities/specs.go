package entities

import (
	"time"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
)

// TODO: Remove annotations when types from service layer are used

type HashicorpSpecs struct {
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

type AkvSpecs struct {
	VaultName    string `json:"vaultName" validate:"required" example:"quorumkeymanager"`
	TenantID     string `json:"tenantID" validate:"required" example:"17255fb0-373b-4a1a-bd47-d211ab86df81"`
	ClientID     string `json:"clientID" validate:"required" example:"8c925036-dd6f-4a1e-a315-5e6fab4f2f09"`
	ClientSecret string `json:"clientSecret" validate:"required" example:"my-secret"`
}

type AwsSpecs struct {
	Region    string `json:"region" validate:"required" example:"eu-west-3"`
	AccessID  string `json:"accessID" validate:"required" example:"AKIAQX7AV2NLJTF5ZZOB"`
	SecretKey string `json:"secretKey" validate:"required" example:"my-secert"`
	Debug     bool   `json:"debug" validate:"required" example:"true"`
}

type LocalKeySpecs struct {
	SecretStore manifest.VaultType `json:"secretStore" validate:"required" example:"HashicorpSecrets"`
	Specs       interface{}        `json:"specs" validate:"required"`
}

type LocalEthSpecs struct {
	Keystore manifest.VaultType `json:"keyStore" validate:"required" example:"HashicorpKeys"`
	Specs    interface{}        `json:"specs" validate:"required"`
}
