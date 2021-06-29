package client

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

const (
	aliasPrefix = "alias/"
)

type AwsKmsClient struct {
	client *kms.KMS
	cfg    *Config
}

func NewKmsClient(cfg *Config) (*AwsKmsClient, error) {
	newSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &AwsKmsClient{
		client: kms.New(newSession, cfg.ToAWSConfig()),
		cfg:    cfg,
	}, nil
}

func (c *AwsKmsClient) CreateKey(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*kms.CreateKeyOutput, error) {
	// Always create with same usage for key now (sign & verify)
	keyUsage := kms.KeyUsageTypeSignVerify

	keySpec, err := c.convertToAWSKeyType(alg)
	if err != nil {
		return nil, err
	}

	out, err := c.client.CreateKey(&kms.CreateKeyInput{
		CustomerMasterKeySpec: &keySpec,
		KeyUsage:              &keyUsage,
		Tags:                  toAWSTags(attr.Tags),
	})
	if err != nil {
		return nil, err
	}

	aliasName := aliasPrefix + id
	_, err = c.client.CreateAlias(&kms.CreateAliasInput{
		AliasName:   &aliasName,
		TargetKeyId: out.KeyMetadata.KeyId,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) GetPublicKey(ctx context.Context, id string) (*kms.GetPublicKeyOutput, error) {
	out, err := c.client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &id,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) DescribeKey(ctx context.Context, id string) (*kms.DescribeKeyOutput, error) {
	out, err := c.client.DescribeKey(&kms.DescribeKeyInput{KeyId: &id})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) ListKeys(ctx context.Context, limit int64, marker string) (*kms.ListKeysOutput, error) {
	input := &kms.ListKeysInput{}
	if limit > 0 {
		input.Limit = &limit
	}
	if len(marker) > 0 {
		input.Marker = &marker
	}

	keys, err := c.client.ListKeys(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return keys, nil
}

func (c *AwsKmsClient) ListTags(ctx context.Context, id, marker string) (*kms.ListResourceTagsOutput, error) {
	input := &kms.ListResourceTagsInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}

	tags, err := c.client.ListResourceTags(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return tags, nil
}

func (c *AwsKmsClient) ListAliases(ctx context.Context, id, marker string) (*kms.ListAliasesOutput, error) {
	input := &kms.ListAliasesInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}

	out, err := c.client.ListAliases(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) Sign(ctx context.Context, id string, msg []byte, signingAlgorithm string) (*kms.SignOutput, error) {
	// Message type is always digest
	msgType := kms.MessageTypeDigest
	out, err := c.client.Sign(&kms.SignInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		SigningAlgorithm: &signingAlgorithm,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) Verify(ctx context.Context, id string, msg, signature []byte, signingAlgorithm string) (*kms.VerifyOutput, error) {
	msgType := kms.MessageTypeDigest
	out, err := c.client.Verify(&kms.VerifyInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		Signature:        signature,
		SigningAlgorithm: &signingAlgorithm,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) DeleteKey(ctx context.Context, id string) (*kms.ScheduleKeyDeletionOutput, error) {
	out, err := c.client.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
		KeyId: &id,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) RecoverKey(ctx context.Context, id string) (*kms.CancelKeyDeletionOutput, error) {
	out, err := c.client.CancelKeyDeletion(&kms.CancelKeyDeletionInput{
		KeyId: &id,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AwsKmsClient) UpdateKey(ctx context.Context, id string, tags map[string]string) (*kms.TagResourceOutput, error) {
	outTagResource, err := c.client.TagResource(&kms.TagResourceInput{
		KeyId: &id,
		Tags:  toAWSTags(tags),
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return outTagResource, nil
}

func (c *AwsKmsClient) convertToAWSKeyType(alg *entities.Algorithm) (string, error) {
	switch {
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		return kms.CustomerMasterKeySpecEccSecgP256k1, nil
	case alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Bn254:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid key type")
	}
}

func toAWSTags(tags map[string]string) []*kms.Tag {
	// populate tags
	var keyTags []*kms.Tag
	for key, value := range tags {
		k, v := key, value
		keyTag := kms.Tag{TagKey: &k, TagValue: &v}
		keyTags = append(keyTags, &keyTag)
	}
	return keyTags
}
