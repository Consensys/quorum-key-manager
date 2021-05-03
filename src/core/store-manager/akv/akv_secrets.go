package akv

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/akv"
)

// Specs is the specs format for an Azure Key Vault secret store
type Specs struct {
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

func NewSecretStore(spec *Specs) (*akv.Store, error) {
	cfg := client.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	store := akv.New(cli)
	return store, nil
}
