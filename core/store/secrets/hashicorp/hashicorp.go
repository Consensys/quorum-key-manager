package hashicorpsecrets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models"
	"io/ioutil"
	"path"
	"strings"
	"time"

	hashicorp "github.com/hashicorp/vault/api"
)

const (
	valueLabel          = "value"
	expirationDateLabel = "expirationDate"
	tagsLabel           = "tags"
	enabledLabel        = "enabled"
)

// Store is an implementation of secret store relying on HashiCorp Vault kv-v2 secret engine
type hashicorpSecretStore struct {
	client *hashicorp.Logical
	cfg    *Config
}

// New creates an HashiCorp secret store
func New(client *hashicorp.Client, cfg *Config) (*hashicorpSecretStore, error) {
	err := client.SetAddress(cfg.Address)
	if err != nil {
		return nil, err
	}

	client.SetNamespace(cfg.Namespace)

	var decodedToken string
	if cfg.TokenFilePath != "" {
		encodedToken, err := ioutil.ReadFile(cfg.TokenFilePath)
		if err != nil {
			return nil, err
		}

		decodedToken = strings.TrimSuffix(string(encodedToken), "\n") // Remove the newline if it exists
		decodedToken = strings.TrimSuffix(decodedToken, "\r")         // This one is for windows compatibility
	} else {
		decodedToken = cfg.Token
	}

	client.SetToken(decodedToken)

	return &hashicorpSecretStore{
		client: client.Logical(),
		cfg:    cfg,
	}, nil
}

// path compute path from hashicorp mount
func (s *hashicorpSecretStore) path(id string) string {
	return path.Join(s.cfg.MountPoint, "data", id)
}

// Set a secret
func (s *hashicorpSecretStore) Set(ctx context.Context, id string, value string, attr *models.Attributes) (*models.Secret, error) {
	data := map[string]interface{}{
		valueLabel:          value,
		expirationDateLabel: attr.ExpireAt.UTC().Format(time.UnixDate),
		tagsLabel:           attr.Tags,
		enabledLabel:        attr.Enabled,
	}

	secret, err := s.client.Write(s.path(id), data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(secret), err
}

// Get a secret
func (s *hashicorpSecretStore) Get(ctx context.Context, id string, version int) (*models.Secret, error) {
	data := map[string][]string{}

	secret, err := s.client.ReadWithData(s.path(id), data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(secret), err
}

// Get all secret ids
func (s *hashicorpSecretStore) List(ctx context.Context) ([]string, error) {
	res, err := s.client.List(path.Join(s.cfg.MountPoint, "metadata"))
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []string{}, nil
	}

	secrets := res.Data["keys"].([]interface{})
	ids := make([]string, len(secrets))
	for i, elem := range secrets {
		ids[i] = fmt.Sprintf("%v", elem)
	}

	return ids, nil
}

// Update a secret
func (s *hashicorpSecretStore) Update(ctx context.Context, id string, newValue string, attr *models.Attributes) (*models.Secret, error) {
	// Update simply overrides a secret
	return s.Set(ctx, id, newValue, attr)
}

func formatHashicorpSecret(secret *hashicorp.Secret) *models.Secret {
	for k, _ := range secret.Data {
		fmt.Println(json.MarshalIndent(secret.Data[k], "", "  "))
	}

	return &models.Secret{
		Value:    secret.Data[valueLabel].(string),
		Attr:     nil,
		Metadata: nil,
	}
}
