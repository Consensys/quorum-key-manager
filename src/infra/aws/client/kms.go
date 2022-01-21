package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/cenkalti/backoff"
)

func (c *AWSClient) CreateKey(ctx context.Context, id, keyType string, tags []*kms.Tag) (*kms.CreateKeyOutput, error) {
	// Always create with same usage for key now (sign & verify)
	keyUsage := kms.KeyUsageTypeSignVerify

	out, err := c.kmsClient.CreateKey(&kms.CreateKeyInput{
		CustomerMasterKeySpec: &keyType,
		KeyUsage:              &keyUsage,
		Tags:                  tags,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	_, err = c.kmsClient.CreateAlias(&kms.CreateAliasInput{
		AliasName:   &id,
		TargetKeyId: out.KeyMetadata.KeyId,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	err = backoff.RetryNotify(func() error {
		descData, err := c.DescribeKey(ctx, id)
		if err != nil {
			return err
		}
		if *descData.KeyMetadata.Enabled {
			return nil
		}
		return errors.New(fmt.Sprintf("keyId %s is still not enabled"))
	}, c.backOff,
		func(err error, t time.Duration) {
			c.logger.Debug(fmt.Sprintf("ERR: %s, retrying in %s", err.Error(), t.String()))
		},
	)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) GetPublicKey(_ context.Context, keyID string) (*kms.GetPublicKeyOutput, error) {
	out, err := c.kmsClient.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &keyID,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) DescribeKey(_ context.Context, id string) (*kms.DescribeKeyOutput, error) {
	out, err := c.kmsClient.DescribeKey(&kms.DescribeKeyInput{KeyId: &id})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) ListKeys(_ context.Context, limit int64, marker string) (*kms.ListKeysOutput, error) {
	input := &kms.ListKeysInput{}
	if limit > 0 {
		input.Limit = &limit
	}
	if len(marker) > 0 {
		input.Marker = &marker
	}

	outListKeys, err := c.kmsClient.ListKeys(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return outListKeys, nil
}

func (c *AWSClient) GetAlias(_ context.Context, keyID string) (string, error) {
	out, err := c.kmsClient.ListAliases(&kms.ListAliasesInput{
		KeyId: aws.String(keyID),
		Limit: aws.Int64(1),
	})
	if err != nil {
		return "", parseKmsErrorResponse(err)
	}

	if len(out.Aliases) > 0 {
		ss := strings.Split(*out.Aliases[0].AliasName, "/")
		return ss[len(ss)-1], nil
	}

	return "", nil
}

func (c *AWSClient) ListTags(_ context.Context, id, marker string) (*kms.ListResourceTagsOutput, error) {
	input := &kms.ListResourceTagsInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}

	tags, err := c.kmsClient.ListResourceTags(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return tags, nil
}

func (c *AWSClient) Sign(_ context.Context, keyID string, msg []byte, signingAlgorithm string) (*kms.SignOutput, error) {
	// Message type is always digest
	msgType := kms.MessageTypeDigest
	out, err := c.kmsClient.Sign(&kms.SignInput{
		KeyId:            &keyID,
		Message:          msg,
		MessageType:      &msgType,
		SigningAlgorithm: &signingAlgorithm,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) DeleteKey(_ context.Context, keyID string) (*kms.ScheduleKeyDeletionOutput, error) {
	out, err := c.kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
		KeyId: &keyID,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) RestoreKey(_ context.Context, keyID string) (*kms.CancelKeyDeletionOutput, error) {
	out, err := c.kmsClient.CancelKeyDeletion(&kms.CancelKeyDeletionInput{
		KeyId: &keyID,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return out, nil
}

func (c *AWSClient) TagResource(_ context.Context, keyID string, tags []*kms.Tag) (*kms.TagResourceOutput, error) {
	outTagResource, err := c.kmsClient.TagResource(&kms.TagResourceInput{
		KeyId: &keyID,
		Tags:  tags,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return outTagResource, nil
}

func (c *AWSClient) UntagResource(_ context.Context, keyID string, tagKeys []*string) (*kms.UntagResourceOutput, error) {
	outUntagResource, err := c.kmsClient.UntagResource(&kms.UntagResourceInput{
		KeyId:   &keyID,
		TagKeys: tagKeys,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}

	return outUntagResource, nil
}
