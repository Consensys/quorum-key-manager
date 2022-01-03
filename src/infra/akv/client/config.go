package client

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/consensys/quorum-key-manager/src/entities"
)

type Config struct {
	Endpoint            string
	SubscriptionID      string
	TenantID            string
	AuxiliaryTenantIDs  string
	ClientID            string
	ClientSecret        string
	CertificatePath     string
	CertificatePassword string
	Username            string
	Password            string
	EnvironmentName     string
	Resource            string
}

func NewConfig(cfg *entities.AzureConfig) *Config {
	return &Config{
		Endpoint:     fmt.Sprintf("https://%s.%s", cfg.VaultName, azure.PublicCloud.KeyVaultDNSSuffix),
		TenantID:     cfg.TenantID,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	}
}

// ToAzureAuthConfig  Inspired by NewAuthorizerFromEnvironmentWithResource from github.com/azure/go-autorest/autorest/azure/auth@v0.5.7/auth.go (https://github.com/Azure/go-autorest/blob/master/autorest/azure/auth/auth.go)
func (c *Config) ToAzureAuthConfig() (autorest.Authorizer, error) {
	resource, err := c.getResource()
	if err != nil {
		return nil, err
	}

	settings, err := c.getSettings()
	if err != nil {
		return nil, err
	}

	settings.Values[auth.Resource] = resource
	return settings.GetAuthorizer()
}

// Inspired by getResource from services/keyvault/auth/auth.go (https://github.com/Azure/azure-sdk-for-go/blob/master/services/keyvault/auth/auth.go)
func (c *Config) getResource() (string, error) {
	var env azure.Environment

	if envName := c.EnvironmentName; envName == "" {
		env = azure.PublicCloud
	} else {
		var err error
		env, err = azure.EnvironmentFromName(envName)
		if err != nil {
			return "", err
		}
	}

	resource := c.Resource
	if resource == "" {
		resource = env.ResourceIdentifiers.KeyVault
	}

	return resource, nil
}

// Inspired  by GetSettingsFromEnvironment from github.com/azure/go-autorest/autorest/azure/auth@v0.5.7/auth.go (https://github.com/Azure/go-autorest/blob/master/autorest/azure/auth/auth.go)
func (c *Config) getSettings() (s auth.EnvironmentSettings, err error) {
	s = auth.EnvironmentSettings{
		Values: map[string]string{},
	}

	if c.SubscriptionID != "" {
		s.Values[auth.SubscriptionID] = c.SubscriptionID
	}
	if c.TenantID != "" {
		s.Values[auth.TenantID] = c.TenantID
	}
	if c.AuxiliaryTenantIDs != "" {
		s.Values[auth.AuxiliaryTenantIDs] = c.AuxiliaryTenantIDs
	}
	if c.ClientID != "" {
		s.Values[auth.ClientID] = c.ClientID
	}
	if c.ClientSecret != "" {
		s.Values[auth.ClientSecret] = c.ClientSecret
	}
	if c.CertificatePath != "" {
		s.Values[auth.CertificatePath] = c.CertificatePath
	}
	if c.CertificatePassword != "" {
		s.Values[auth.CertificatePassword] = c.CertificatePassword
	}
	if c.Username != "" {
		s.Values[auth.Username] = c.Username
	}
	if c.Password != "" {
		s.Values[auth.Password] = c.Password
	}
	if c.EnvironmentName != "" {
		s.Values[auth.EnvironmentName] = c.EnvironmentName
	}
	if c.Resource != "" {
		s.Values[auth.Resource] = c.Resource
	}

	if v := s.Values[auth.EnvironmentName]; v == "" {
		s.Environment = azure.PublicCloud
	} else {
		s.Environment, err = azure.EnvironmentFromName(v)
	}
	if s.Values[auth.Resource] == "" {
		s.Values[auth.Resource] = s.Environment.ResourceManagerEndpoint
	}

	return
}
