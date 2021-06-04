package aws

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/infra/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
)

// Store is an implementation of key store relying on AWS kms
type KeyStore struct {
	client aws.KmsClient
	logger *log.Logger
}

// New creates an AWS secret store
func New(client aws.KmsClient, logger *log.Logger) *KeyStore {
	return &KeyStore{
		client: client,
		logger: logger,
	}
}

// Info returns store information
func (ks *KeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

// Create a new key and stores it
func (ks *KeyStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := ks.logger.WithField("id", id)

	key, err := ks.client.CreateKey(ctx, id, alg, attr)
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, err
	}

	publicKeyOut, err := ks.client.GetPublicKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve pub key info")
		return nil, err
	}

	algo := algoFromAWSPublicKeyInfo(publicKeyOut)

	outKey := &entities.Key{
		ID:        *key.KeyMetadata.KeyId,
		PublicKey: publicKeyOut.PublicKey,
		Algo:      algo,
		Metadata: &entities.Metadata{

			Disabled:  !*key.KeyMetadata.Enabled,
			CreatedAt: *key.KeyMetadata.CreationDate,
			DeletedAt: *key.KeyMetadata.DeletionDate,
		},
		Tags: nil,
	}
	logger.Info("key created successfully")
	return outKey, nil
}

// Import an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (ks *KeyStore) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

// Get the public part of a stored key.
func (ks *KeyStore) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := ks.logger.WithField("id", id)
	outGetKey, err := ks.client.GetPublicKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get public key")
		return nil, err
	}

	retKey := &entities.Key{ID: *outGetKey.KeyId,
		PublicKey: outGetKey.PublicKey}

	// List associated tags
	tags := make(map[string]string)
	nextMarker := ""
	for {
		ret, err := ks.client.ListTags(ctx, id, nextMarker)
		if err != nil {
			logger.WithError(err).Error("failed to list key tags")
			return nil, err
		}

		for _, tag := range ret.Tags {
			tags[*tag.TagKey] = *tag.TagValue
		}
		if !*ret.Truncated {
			break
		}
		nextMarker = *ret.NextMarker
	}
	logger.Info("successfully got key info")
	return retKey, nil
}

// List keys
func (ks *KeyStore) List(ctx context.Context) ([]string, error) {
	var keys []string
	nextMarker := ""

	// Loop until the entire list is constituted
	for {
		ret, err := ks.client.ListKeys(ctx, 0, nextMarker)
		if err != nil {
			ks.logger.WithError(err).Error("failed to list keys")
			return nil, err
		}

		for _, key := range ret.Keys {
			keys = append(keys, *key.KeyId)
		}

		if ret.NextMarker == nil {
			break
		}
		nextMarker = *ret.NextMarker

	}
	ks.logger.Info("keys listed successfully")
	return keys, nil
}

// Update key tags
func (ks *KeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

// Delete key not permanently, by using Undelete() the key can be enabled again
func (ks *KeyStore) Delete(ctx context.Context, id string) error {
	logger := ks.logger.WithField("id", id)
	_, err := ks.client.DeleteKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return err
	}
	logger.Info("deleted key successfully")
	return err
}

// GetDeleted keys
func (ks *KeyStore) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

// ListDeleted keys
func (ks *KeyStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

// Undelete a previously deleted secret
func (ks *KeyStore) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Destroy secret permanently
func (ks *KeyStore) Destroy(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Sign from any arbitrary data using the specified key
func (ks *KeyStore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := ks.logger.WithField("id", id)
	outSignature, err := ks.client.Sign(ctx, id, data)
	if err != nil {
		logger.WithError(err).Error("failed to sign")
		return nil, err
	}
	logger.Info("data signed successfully")
	return outSignature.Signature, nil
}

// Encrypt any arbitrary data using a specified key
func (ks *KeyStore) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

// Decrypt a single block of encrypted data.
func (ks *KeyStore) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
