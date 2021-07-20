package akv

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/akv/client"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
)

// SecretSpecs is the specs format for an Azure Key Vault secret store
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

func NewSecretStore(spec *SecretSpecs, logger log.Logger) (*akv.Store, error) {
	cfg := client.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client.NewClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AKV client (secrets)"
		logger.WithError(err).Error(errMessage, "specs", spec)
		return nil, errors.ConfigError(errMessage)
	}

	store := akv.New(cli, logger)
	return store, nil
}
