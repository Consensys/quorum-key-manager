package entities

import (
	"time"
)

const (
	HashicorpVaultType = "hashicorp"
	AzureVaultType     = "azure"
	AWSVaultType       = "aws"
)

type Vault struct {
	Client    interface{}
	VaultType string
	Name      string
}

type HashicorpConfig struct {
	MountPoint    string        `json:"mountPoint" yaml:"mount_point" validate:"required" example:"secret"`
	Address       string        `json:"address"  yaml:"address" validate:"required" example:"https://hashicorp:8200"`
	Token         string        `json:"token,omitempty" yaml:"token" example:"s.W7IMlFuBGsTaR6uHLcGDw9Mq"`
	TokenPath     string        `json:"tokenPath,omitempty" yaml:"token_path,omitempty" example:"/vault/token/.my_token"`
	Namespace     string        `json:"namespace,omitempty" yaml:"namespace,omitempty" example:"namespace"`
	CACert        string        `json:"CACert,omitempty" yaml:"ca_cert,omitempty" example:"/vault/tls/ca.crt"`
	CAPath        string        `json:"CAPath,omitempty" yaml:"ca_path,omitempty" example:"/vault/tls/root-certs"`
	ClientCert    string        `json:"clientCert,omitempty" yaml:"client_cert,omitempty" example:"/vault/tls/client.crt"`
	ClientKey     string        `json:"clientKey,omitempty" yaml:"client_key,omitempty" example:"/vault/tls/client.key"`
	TLSServerName string        `json:"TLSServerName,omitempty" yaml:"tls_server_name,omitempty" example:"server-name"`
	ClientTimeout time.Duration `json:"clientTimeout,omitempty" yaml:"client_timeout,omitempty" example:"60s"`
	RateLimit     float64       `json:"rateLimit,omitempty" yaml:"rate_limit,omitempty" example:"0"`
	BurstLimit    int           `json:"burstLimit,omitempty" yaml:"burst_limit,omitempty" example:"0"`
	MaxRetries    int           `json:"maxRetries,omitempty" yaml:"max_retries,omitempty" example:"2"`
	SkipVerify    bool          `json:"skipVerify,omitempty" yaml:"skip_verify,omitempty" example:"false"`
}

type AzureConfig struct {
	VaultName    string `json:"vaultName" yaml:"vault_name" validate:"required" example:"quorumkeymanager"`
	TenantID     string `json:"tenantID" yaml:"tenant_id" validate:"required" example:"17255fb0-373b-4a1a-bd47-d211ab86df81"`
	ClientID     string `json:"clientID" yaml:"client_id" validate:"required" example:"8c925036-dd6f-4a1e-a315-5e6fab4f2f09"`
	ClientSecret string `json:"clientSecret" yaml:"client_secret" validate:"required" example:"my-secret"`
}

type AWSConfig struct {
	Region    string `json:"region" yaml:"region" validate:"required" example:"eu-west-3"`
	AccessID  string `json:"accessID" yaml:"access_id" validate:"required" example:"AKIAQX7AV2NLJTF5ZZOB"`
	SecretKey string `json:"secretKey" yaml:"secret_key" validate:"required" example:"my-secert"`
	Debug     bool   `json:"debug" yaml:"debug" validate:"required" example:"true"`
}
