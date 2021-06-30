package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

const (
	aliasPrefix = "alias"
)

type KeyStore struct {
	client aws.KmsClient
	logger log.Logger
}

func New(client aws.KmsClient, logger log.Logger) *KeyStore {
	return &KeyStore{
		client: client,
		logger: logger,
	}
}

func (s *KeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *KeyStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)
	logger.Debug("creating key")

	_, err := s.client.CreateKey(ctx, alias(id), alg, attr)
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, err
	}

	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info("key created successfully")
	return key, nil
}

// Import an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (s *KeyStore) Import(_ context.Context, _ string, _ []byte, _ *entities.Algorithm, _ *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	outDescribe, err := s.client.DescribeKey(ctx, alias(id))
	if err != nil {
		logger.WithError(err).Error("failed to describe key")
		return nil, err
	}
	keyID := *outDescribe.KeyMetadata.KeyId

	outGetKey, err := s.client.GetPublicKey(ctx, keyID)
	if err != nil {
		logger.WithError(err).Error("failed to get public key")
		return nil, err
	}

	tags, err := s.listTags(ctx, keyID)
	if err != nil {
		logger.WithError(err).Error("failed to list tags")
		return nil, err
	}

	logger.Debug("key retrieved successfully")
	return &entities.Key{
		ID:          id,
		PublicKey:   outGetKey.PublicKey,
		Algo:        parseAlgorithm(outGetKey),
		Metadata:    parseMetadata(outDescribe),
		Tags:        tags,
		Annotations: parseAnnotations(*outGetKey.KeyId, outDescribe),
	}, nil
}

func (s *KeyStore) List(ctx context.Context) ([]string, error) {
	keyIDs, err := s.List(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to list keys")
		return nil, err
	}

	var ids []string
	for _, keyID := range keyIDs {
		outKey, err := s.client.DescribeKey(ctx, keyID)
		if err != nil {
			s.logger.WithError(err).Error("failed to get key for alias listing")
			return nil, err
		}

		keyAlias, err := s.client.GetAlias(ctx, *outKey.KeyMetadata.KeyId)
		if err != nil {
			s.logger.WithError(err).Error("failed to get key for alias listing")
			return nil, err
		}

		ids = append(ids, keyAlias)
	}

	s.logger.Debug("keys listed successfully")
	return ids, nil
}

func (s *KeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)
	logger.Debug("updating key")

	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.client.UpdateKey(ctx, key.Annotations[awsKeyID], attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to update key")
		return nil, err
	}

	logger.Info("key updated successfully")
	return s.Get(ctx, id)
}

func (s *KeyStore) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("deleting key")

	key, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteKey(ctx, key.Annotations[awsKeyID])
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}

func (s *KeyStore) GetDeleted(_ context.Context, _ string) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) ListDeleted(_ context.Context) ([]string, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) Undelete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("restoring key")

	_, err := s.client.RestoreKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to restore key")
		return err
	}

	logger.Info("key restored successfully")
	return nil
}

// Destroy destroys an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (s *KeyStore) Destroy(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *KeyStore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := s.logger.With("id", id)
	logger.Debug("signing payload")

	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Only sign with ECDSA, extract the algorithm from the key when more keys are available
	signature, err := s.client.Sign(ctx, key.Annotations[awsKeyID], data, kms.SigningAlgorithmSpecEcdsaSha256)
	if err != nil {
		logger.WithError(err).Error("failed to sign")
		return nil, err
	}

	logger.Debug("payload signed successfully")
	return signature.Signature, nil
}

func (s *KeyStore) Verify(ctx context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	return keys.VerifySignature(s.logger, pubKey, data, sig, algo)
}

func (s *KeyStore) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *KeyStore) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *KeyStore) listTags(ctx context.Context, keyID string) (map[string]string, error) {
	tags := make(map[string]string)

	nextMarker := ""
	for {
		ret, err := s.client.ListTags(ctx, keyID, nextMarker)
		if err != nil {
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

	return tags, nil
}

func alias(id string) string {
	return fmt.Sprintf("%s/%s", aliasPrefix, id)
}
