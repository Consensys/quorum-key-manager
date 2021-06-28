package aws

import (
	"context"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

const (
	awskeyIDTag           = "aws-KeyID"
	awsCustomerKeyStoreID = "aws-KeyStoreID"
	awsCloudHsmClusterID  = "aws-ClusterHSMID"
	awsAccountID          = "aws-AccountID"
	awsARN                = "awsARN"
)

// Store is an implementation of key store relying on AWS kms
type KeyStore struct {
	client aws.KmsClient
	logger log.Logger
}

// New creates an AWS secret store
func New(client aws.KmsClient, logger log.Logger) *KeyStore {
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
	logger := ks.logger

	key, alias, err := ks.client.CreateKey(ctx, id, alg, attr)
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, err
	}

	publicKeyOut, err := ks.client.GetPublicKey(ctx, *key.KeyMetadata.KeyId)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve pub key info")
		return nil, err
	}

	algo := algoFromAWSPublicKeyInfo(publicKeyOut)

	// List associated tags
	tags, err2 := ks.doListTags(ctx, *key.KeyMetadata.KeyId, logger)
	if err2 != nil {
		return nil, err2
	}

	outKey := &entities.Key{
		ID:        *alias,
		PublicKey: publicKeyOut.PublicKey,
		Algo:      algo,
		Metadata:  metadataFromAWSKey(key),
		Tags:      tags,
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
	logger := ks.logger.With("id", id)
	outGetKey, err := ks.client.GetPublicKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get public key")
		return nil, err
	}

	retKey := &entities.Key{
		PublicKey: outGetKey.PublicKey,
		Algo:      algoFromAWSPublicKeyInfo(outGetKey),
	}

	// List associated tags
	tags, err2 := ks.doListTags(ctx, *outGetKey.KeyId, logger)
	if err2 != nil {
		return nil, err2
	}
	// List aliases

	aliasMarker := ""
	_, cleanAliases, errListAliases := ks.client.ListAliases(ctx, *outGetKey.KeyId, aliasMarker)
	if errListAliases != nil {
		logger.WithError(errListAliases).Error("failed to list key aliases")
		return nil, errListAliases
	}

	for _, alias := range cleanAliases {
		// KeyID will be first aliasName found
		retKey.ID = alias
		break
	}
	// Describe key
	retDescribe, err := ks.client.DescribeKey(ctx, *outGetKey.KeyId)
	if err != nil {
		logger.WithError(err).Error("failed to describe key")
		return nil, err
	}

	// Populate associated annotations
	annotations := make(map[string]string)
	// First set aws ID
	annotations[awskeyIDTag] = *outGetKey.KeyId

	fillAwsAnnotations(annotations, retDescribe)

	retKey.Tags = tags
	retKey.Annotations = annotations

	// Populate metadata
	retKey.Metadata = metadataFromAWSDescribeKey(retDescribe)

	logger.Info("successfully got key info")
	return retKey, nil
}

func (ks *KeyStore) doListTags(ctx context.Context, keyID string, logger log.Logger) (map[string]string, error) {
	tags := make(map[string]string)

	nextMarker := ""
	for {
		ret, errListTags := ks.client.ListTags(ctx, keyID, nextMarker)
		if errListTags != nil {

			logger.WithError(errListTags).Error("failed to list key tags")
			return nil, errListTags
		}

		for _, tag := range ret.Tags {
			tags[*tag.TagKey] = *tag.TagValue
		}
		if !*ret.Truncated {
			break
		}
		nextMarker = *ret.NextMarker
	}
	return tags, nil
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
	logger := ks.logger.With("id", id)

	_, err := ks.client.UpdateKey(ctx, id, attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to update key")
		return nil, err
	}

	logger.Info("updated key successfully")

	return ks.Get(ctx, id)
}

// Delete key not permanently, by using Undelete() the key can be enabled again
func (ks *KeyStore) Delete(ctx context.Context, id string) error {
	logger := ks.logger.With("id", id)
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
	logger := ks.logger.With("id", id)
	outSignature, err := ks.client.Sign(ctx, id, data)
	if err != nil {
		logger.WithError(err).Error("failed to sign")
		return nil, err
	}
	logger.Info("data signed successfully")
	return outSignature.Signature, nil
}

func (ks *KeyStore) Verify(ctx context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	return keys.VerifySignature(ks.logger, pubKey, data, sig, algo)
}

// Encrypt any arbitrary data using a specified key
func (ks *KeyStore) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

// Decrypt a single block of encrypted data.
func (ks *KeyStore) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
