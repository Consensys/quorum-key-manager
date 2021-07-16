package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
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
	var keyType string

	switch {
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		keyType = kms.CustomerMasterKeySpecEccSecgP256k1
	default:
		errMessage := "invalid or not supported elliptic curve and signing algorithm for AWS key creation"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
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
func (s *KeyStore) Import(_ context.Context, _ string, _ []byte, _ *entities.Algorithm, _ *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) Get(ctx context.Context, id string) (*entities.Key, error) {
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

func (s *KeyStore) List(ctx context.Context) ([]string, error) {
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

func (s *KeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	keyID := key.Annotations[awsKeyID]
	logger := s.logger.With("id", id, "key_id", keyID)

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

func (s *KeyStore) Delete(ctx context.Context, id string) error {
	key, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteKey(ctx, key.Annotations[awsKeyID])
	if err != nil {
		errMessage := "failed to delete AWS key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *KeyStore) GetDeleted(_ context.Context, _ string) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) ListDeleted(_ context.Context) ([]string, error) {
	return nil, errors.ErrNotSupported
}

func (s *KeyStore) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Destroy destroys an externally created key and stores it
// this feature is not supported by AWS kms
// always returns errors.ErrNotSupported
func (s *KeyStore) Destroy(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *KeyStore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Only sign with ECDSA, extract the algorithm from the key when more keys are available
	outSignature, err := s.client.Sign(ctx, key.Annotations[awsKeyID], data, kms.SigningAlgorithmSpecEcdsaSha256)
	if err != nil {
		errMessage := "failed to sign using AWS key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	signature, err := parseSignature(outSignature)
	if err != nil {
		errMessage := "failed to parse signature from AWS"
		s.logger.With("id", id, "signature", signature).WithError(err).Error(errMessage)
		return nil, errors.AWSError(errMessage)
	}

	return signature, nil
}

func (s *KeyStore) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
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
