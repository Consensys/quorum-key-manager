package hashicorpsecrets

import (
	"context"
	"fmt"
	"path"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
	"github.com/hashicorp/vault/api"
)

// Store is an implementation of secret store relying on hashicorp vault
type Store struct {
	hashicorp *api.Client
	cfg       *Config
}

// New creates an hasicorp secret store
func New(cfg *Config) (*Store, error) {
	// Create Hashicorp client
	hashicorp, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	// Set client config from cfg
	err = hashicorp.SetAddress(cfg.Addr)
	if err != nil {
		return nil, err
	}

	hashicorp.SetNamespace(cfg.Namespace)
	hashicorp.SetToken(cfg.Token)

	// Return store
	return &Store{
		hashicorp: hashicorp,
		cfg:       cfg,
	}, nil
}

// path compute path from hashicorp mount
func (s *Store) path(id string) string {
	return path.Join(s.cfg.Mount, id)
}

// Set secret
func (s *Store) Set(ctx context.Context, id string, value []byte, attr *types.Attributes) (*types.Secret, error) {
	data := map[string]interface{}{
		// TODO: compute hashicorp data
		"value": string(value),
	}

	// Set tags
	for k, v := range attr.Tags {
		data[k] = v
	}

	secret, err := s.hashicorp.Logical().Write(s.path(id), data)
	if err != nil {
		return nil, err
	}

	return FromHashicorpSecret(secret), err
}

// Get a secret
func (s *Store) Get(ctx context.Context, id string, version int) (*types.Secret, error) {
	data := map[string][]string{
		// TODO: compute hashicorp data
		"version": []string{fmt.Sprintf("%v", version)},
	}

	secret, err := s.hashicorp.Logical().ReadWithData(s.path(id), data)
	if err != nil {
		return nil, err
	}

	return FromHashicorpSecret(secret), err
}

func FromHashicorpSecret(secret *api.Secret) *types.Secret {
	return &types.Secret{
		// TODO: compute secret from hashicorp api.Secret
	}
}
