package client

import (
	"context"
	"os"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/aws/aws-sdk-go/service/kms"
)

const (
	aliasPrefix = "alias/"
)

func (c *AwsKmsClient) CreateKey(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*kms.CreateKeyOutput, *string, error) {
	// Always create with same usage for key now (sign & verify)
	keyUsage := kms.KeyUsageTypeSignVerify

	keySpec, err := convertToAWSKeyType(alg)
	if err != nil {
		return nil, nil, err
	}

	outKey, err := c.client.CreateKey(&kms.CreateKeyInput{
		CustomerMasterKeySpec: &keySpec,
		KeyUsage:              &keyUsage,
		Tags:                  toAWSTags(attr),
	})

	if err != nil {
		return nil, nil, err
	}

	aliasName := aliasPrefix + id
	_, err = c.client.CreateAlias(&kms.CreateAliasInput{
		AliasName:   &aliasName,
		TargetKeyId: outKey.KeyMetadata.KeyId,
	})

	if err != nil {
		return nil, nil, parseKmsErrorResponse(err)
	}

	// Retrieve first alias found and assign it to keyID
	var retAlias *string
	_, aliasList, err := c.ListAliases(ctx, *outKey.KeyMetadata.KeyId, "")

	if len(aliasList) > 0 {
		retAlias = &aliasList[0]
	}
	if err != nil {
		return nil, nil, parseKmsErrorResponse(err)
	}
	return outKey, retAlias, nil
}

func (c *AwsKmsClient) GetPublicKey(ctx context.Context, id string) (*kms.GetPublicKeyOutput, error) {
	outGetPubKey, err := c.client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &id,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outGetPubKey, nil
}

func (c *AwsKmsClient) DescribeKey(ctx context.Context, id string) (*kms.DescribeKeyOutput, error) {
	input := &kms.DescribeKeyInput{KeyId: &id}

	outDescribeKey, err := c.client.DescribeKey(input)

	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outDescribeKey, nil
}

func (c *AwsKmsClient) ListKeys(ctx context.Context, limit int64, marker string) (*kms.ListKeysOutput, error) {
	input := &kms.ListKeysInput{}
	if limit > 0 {
		input.Limit = &limit
	}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	outListKeys, err := c.client.ListKeys(input)
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outListKeys, nil
}

func (c *AwsKmsClient) ListTags(ctx context.Context, id, marker string) (*kms.ListResourceTagsOutput, error) {
	input := &kms.ListResourceTagsInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	outListTags, err := c.client.ListResourceTags(input)

	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outListTags, nil
}

func (c *AwsKmsClient) ListAliases(ctx context.Context, id, marker string) (*kms.ListAliasesOutput, []string, error) {
	input := &kms.ListAliasesInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	outListAliases, err := c.client.ListAliases(input)
	if err != nil {
		return nil, nil, parseKmsErrorResponse(err)
	}

	cleanList := []string{}
	for _, awsalias := range outListAliases.Aliases {
		alias := strings.Replace(*awsalias.AliasName, aliasPrefix, "", -1)
		cleanList = append(cleanList, alias)
	}

	if err != nil {
		return nil, nil, parseKmsErrorResponse(err)
	}
	return outListAliases, cleanList, nil

}

func (c *AwsKmsClient) Sign(ctx context.Context, id string, msg []byte) (*kms.SignOutput, error) {
	// Message type is always digest
	msgType := kms.MessageTypeDigest
	signingAlg := kms.SigningAlgorithmSpecEcdsaSha256
	outSign, err := c.client.Sign(&kms.SignInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		SigningAlgorithm: &signingAlg,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outSign, nil
}

func (c *AwsKmsClient) Verify(ctx context.Context, id string, msg, signature []byte) (*kms.VerifyOutput, error) {
	msgType := kms.MessageTypeDigest
	signingAlg := kms.SigningAlgorithmSpecEcdsaSha256
	outVerify, err := c.client.Verify(&kms.VerifyInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		Signature:        signature,
		SigningAlgorithm: &signingAlg,
	})

	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outVerify, nil
}

func (c *AwsKmsClient) DeleteKey(ctx context.Context, id string) (*kms.DisableKeyOutput, error) {
	outDisable, err := c.client.DisableKey(&kms.DisableKeyInput{
		KeyId: &id,
	})
	if err != nil {
		return nil, parseKmsErrorResponse(err)
	}
	return outDisable, nil
}

func convertToAWSKeyType(alg *entities.Algorithm) (string, error) {
	switch alg.Type {
	case entities.Ecdsa:
		if alg.EllipticCurve == entities.Secp256k1 {
			if isTestOn() {
				return kms.CustomerMasterKeySpecEccNistP256, nil
			}
			return kms.CustomerMasterKeySpecEccSecgP256k1, nil
		}
		return "", errors.InvalidParameterError("invalid curve")
	case entities.Eddsa:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid key type")
	}
}

func toAWSTags(attr *entities.Attributes) []*kms.Tag {
	// populate tags
	var keyTags []*kms.Tag
	for key, value := range attr.Tags {
		k, v := key, value
		keyTag := kms.Tag{TagKey: &k, TagValue: &v}
		keyTags = append(keyTags, &keyTag)
	}
	return keyTags
}

func isTestOn() bool {
	val, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return false
	}
	return strings.EqualFold("test", val)
}
