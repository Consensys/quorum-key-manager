package akv

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	client2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/akv/client"
	akv2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/akv"
)

// Specs is the specs format for an Azure Key Vault key store
type KeySpecs struct {
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

func NewKeyStore(spec *KeySpecs, logger *log.Logger) (*akv2.Store, error) {
	cfg := client2.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client2.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	store := akv2.New(cli, logger)
	return store, nil
}
