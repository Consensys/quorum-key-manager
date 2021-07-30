package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

const (
	aliasPrefix = "alias"
)

type Store struct {
	client aws.KmsClient
	db     database.Keys // AWS key store needs the DB to be able to retrieve the key and get the AWS keyID from Annotations
	logger log.Logger
}

var _ keys.Store = &Store{}

func New(client aws.KmsClient, db database.Keys, logger log.Logger) *Store {
	return &Store{
		client: client,
		db:     db,
		logger: logger,
	}
}

func (s *Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)

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

func (s *Store) Import(_ context.Context, _ string, _ []byte, _ *entities.Algorithm, _ *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	key, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}
	keyID := key.Annotations.AWSKeyID
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

func (s *Store) Delete(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *Store) Undelete(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	key, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteKey(ctx, key.Annotations.AWSKeyID)
	if err != nil {
		errMessage := "failed to permanently delete AWS key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Sign(ctx context.Context, id string, data []byte, _ *entities.Algorithm) ([]byte, error) {
	key, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Only sign with ECDSA, extract the algorithm from the key when more keys are available
	outSignature, err := s.client.Sign(ctx, key.Annotations.AWSKeyID, data, kms.SigningAlgorithmSpecEcdsaSha256)
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

func (s *Store) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
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

func (s *Store) get(ctx context.Context, id string) (*entities.Key, error) {
	return s.db.Get(ctx, id)
}

func alias(id string) string {
	return fmt.Sprintf("%s/%s", aliasPrefix, id)
}
