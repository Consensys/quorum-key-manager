package akv

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/akv/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys/akv"
)

// KeySpecs is the specs format for an Azure Key Vault key store
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

func NewKeyStore(spec *KeySpecs, logger log.Logger) (*akv.Store, error) {
	cfg := client.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client.NewClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AKV client (keys)"
		logger.WithError(err).Error(errMessage, "specs", spec)
		return nil, errors.ConfigError(errMessage)
	}

	store := akv.New(cli, logger)
	return store, nil
}
