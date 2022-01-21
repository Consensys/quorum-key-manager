package aws

import (
	"context"
	"fmt"

	entities2 "github.com/consensys/quorum-key-manager/src/entities"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

const (
	aliasPrefix = "alias"
)

type Store struct {
	client aws.KmsClient
	logger log.Logger
}

var _ stores.KeyStore = &Store{}

func New(client aws.KmsClient, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Info(context.Context) (*entities.Store, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	var keyType string

	switch {
	case alg.Type == entities2.Ecdsa && alg.EllipticCurve == entities2.Secp256k1:
		keyType = kms.CustomerMasterKeySpecEccSecgP256k1
	default:
		errMessage := "invalid or not supported elliptic curve and signing algorithm for AWS key creation"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.NotSupportedError(errMessage)
	}

	_, err := s.client.CreateKey(ctx, alias(id), keyType, toTags(attr.Tags))
	if err != nil {
		errMessage := "failed to create AWS key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Import an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (s *Store) Import(_ context.Context, _ string, _ []byte, _ *entities2.Algorithm, _ *entities.Attributes) (*entities.Key, error) {
	err := errors.NotSupportedError("import secret is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	outDescribe, err := s.client.DescribeKey(ctx, alias(id))
	if err != nil {
		errMessage := "failed to get AWS key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}
	keyID := *outDescribe.KeyMetadata.KeyId
	logger = logger.With("key_id", keyID)

	outPublicKey, err := s.client.GetPublicKey(ctx, keyID)
	if err != nil {
		errMessage := "failed to get AWS public key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	tags, err := s.listTags(ctx, keyID)
	if err != nil {
		errMessage := "failed to list AWS key tags"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	key, err := parseKey(id, outPublicKey, outDescribe, tags)
	if err != nil {
		errMessage := "failed to parse key retrieved from AWS KMS"
		logger.WithError(err).Error(errMessage)
		return nil, errors.AWSError(errMessage)
	}

	return key, nil
}

func (s *Store) List(ctx context.Context, _, _ uint64) ([]string, error) {
	var ids []string
	nextMarker := ""

	// Loop until the entire list is constituted
	for {
		ret, err := s.client.ListKeys(ctx, 0, nextMarker)
		if err != nil {
			errMessage := "failed to list AWS keys"
			s.logger.WithError(err).Error(errMessage)
			return nil, errors.FromError(err).SetMessage(errMessage)
		}

		for _, key := range ret.Keys {
			keyAlias, err := s.client.GetAlias(ctx, *key.KeyId)
			if err != nil {
				errMessage := "failed to get AWS key alias"
				s.logger.With("keyID", *key.KeyId).WithError(err).Error(errMessage)
				return nil, errors.FromError(err).SetMessage(errMessage)
			}

			// We should not crash if not alias is found even if this should never happen is using the QKM
			if keyAlias != "" {
				ids = append(ids, keyAlias)
			}
		}

		if ret.NextMarker == nil {
			break
		}
		nextMarker = *ret.NextMarker
	}

	return ids, nil
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)
	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	keyID := key.Annotations.AWSKeyID

	tagKeys := make([]*string, len(key.Tags))
	i := 0
	for k := range key.Tags {
		tagKey := k
		tagKeys[i] = &tagKey
		i++
	}

	_, err = s.client.UntagResource(ctx, keyID, tagKeys)
	if err != nil {
		errMessage := "failed to untag AWS key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	_, err = s.client.TagResource(ctx, keyID, toTags(attr.Tags))
	if err != nil {
		errMessage := "failed to tag AWS key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	key.Tags = attr.Tags
	return key, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	keyID, err := s.getAWSKeyID(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteKey(ctx, keyID)
	if err != nil {
		errMessage := "failed to delete AWS key"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) GetDeleted(_ context.Context, _ string) (*entities.Key, error) {
	err := errors.NotSupportedError("get deleted key is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) ListDeleted(_ context.Context, _, _ uint64) ([]string, error) {
	err := errors.NotSupportedError("list deleted keys is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) Restore(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	keyID, err := s.getAWSKeyID(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.RestoreKey(ctx, keyID)
	if err != nil {
		errMessage := "failed to restore AWS key"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

// Destroy destroys an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (s *Store) Destroy(_ context.Context, _ string) error {
	err := errors.NotSupportedError("destroy key is not supported")
	s.logger.Warn(err.Error())
	return err
}

func (s *Store) Sign(ctx context.Context, id string, data []byte, _ *entities2.Algorithm) ([]byte, error) {
	logger := s.logger.With("id", id)
	keyID, err := s.getAWSKeyID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Only sign with ECDSA, extract the algorithm from the key when more keys are available
	outSignature, err := s.client.Sign(ctx, keyID, data, kms.SigningAlgorithmSpecEcdsaSha256)
	if err != nil {
		errMessage := "failed to sign using AWS key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	signature, err := parseSignature(outSignature)
	if err != nil {
		errMessage := "failed to parse signature from AWS"
		logger.WithError(err).Error(errMessage)
		return nil, errors.AWSError(errMessage)
	}

	return signature, nil
}

func (s *Store) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities2.Algorithm) error {
	err := errors.NotSupportedError("verify signature is not supported")
	s.logger.Warn(err.Error())
	return err
}

func (s *Store) Encrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) getAWSKeyID(ctx context.Context, id string) (string, error) {
	outDescribe, err := s.client.DescribeKey(ctx, alias(id))
	if err != nil {
		errMessage := "failed to get AWS keyID"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return "", errors.FromError(err).SetMessage(errMessage)
	}

	return *outDescribe.KeyMetadata.KeyId, nil
}

func (s *Store) listTags(ctx context.Context, keyID string) (map[string]string, error) {
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
