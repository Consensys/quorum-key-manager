package akv

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/akv/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets/akv"
)

// Specs is the specs format for an Azure Key Vault secret store
type SecretSpecs struct {
	VaultName           string `json:"vaultName"`
	SubscriptionID      string `json:"subscriptionID"`
	TenantID            string `json:"tenantID"`
	AuxiliaryTenantIDs  string `json:"auxiliaryTenantIDs"`
	ClientID            string `json:"clientID"`
	ClientSecret        string `json:"clientSecret"`
	CertificatePath     string `json:"certificatePath"`
	CertificatePassword string `json:"certificatePassword"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	EnvironmentName     string `json:"environmentName"`
	Resource            string `json:"resource"`
}

func NewSecretStore(spec *SecretSpecs, logger *log_old.Logger) (*akv.Store, error) {
	cfg := client.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	store := akv.New(cli, logger)
	return store, nil
}
