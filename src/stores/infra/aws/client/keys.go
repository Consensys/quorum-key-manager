package client

import (
	"context"
	"os"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
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
		return nil, nil, translateAwsKmsError(err)
	}

	// Get confirmation alias was created
	var aliasCreated bool
	var retAlias *string
	listAliasOutput, err := c.client.ListAliases(&kms.ListAliasesInput{
		KeyId: outKey.KeyMetadata.KeyId,
	})

	if err != nil {
		return nil, nil, translateAwsKmsError(err)
	}

	for _, listedAlias := range listAliasOutput.Aliases {
		if strings.Contains(*listedAlias.AliasName, aliasName) {
			aliasCreated = true
			break
		}
	}

	if aliasCreated {
		retAlias = &aliasName
	}

	return outKey, retAlias, nil
}

func (c *AwsKmsClient) GetPublicKey(ctx context.Context, id string) (*kms.GetPublicKeyOutput, error) {
	outGetPubKey, err := c.client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &id,
	})
	return outGetPubKey, translateAwsKmsError(err)
}

func (c *AwsKmsClient) DescribeKey(ctx context.Context, id string) (*kms.DescribeKeyOutput, error) {
	input := &kms.DescribeKeyInput{KeyId: &id}

	outDescribeKey, err := c.client.DescribeKey(input)

	return outDescribeKey, translateAwsKmsError(err)
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
	return outListKeys, translateAwsKmsError(err)
}

func (c *AwsKmsClient) ListTags(ctx context.Context, id, marker string) (*kms.ListResourceTagsOutput, error) {
	input := &kms.ListResourceTagsInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	outListTags, err := c.client.ListResourceTags(input)

	return outListTags, translateAwsKmsError(err)
}

func (c *AwsKmsClient) ListAliases(ctx context.Context, id, marker string) (*kms.ListAliasesOutput, []string, error) {
	input := &kms.ListAliasesInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	outListAliases, err := c.client.ListAliases(input)
	if err != nil {
		return nil, nil, translateAwsKmsError(err)
	}

	cleanList := []string{}
	for _, awsalias := range outListAliases.Aliases {
		alias := strings.Replace(*awsalias.AliasName, aliasPrefix, "", -1)
		cleanList = append(cleanList, alias)
	}

	return outListAliases, cleanList, nil

}

// ImportKey(ctx context.Context, input *kms.ImportKeyMaterialInput, tags map[string]string) (*kms.ImportKeyMaterialOutput, error)
// ImportKey(ctx context.Context, input *kms.ImportKeyMaterialInput, tags map[string]string) (*kms.ImportKeyMaterialOutput, error)

// GetKey(ctx context.Context, name string, version string) (keyvault.KeyBundle, error)
/*
UpdateKey(ctx context.Context, input *kms.UpdateCustomKeyStoreInput, tags map[string]string) (*kms.UpdateCustomKeyStoreOutput, error)
DeleteKey(ctx context.Context, keyName string) (*kms.DeleteCustomKeyStoreOutput, error)
GetDeletedKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error)
GetDeletedKeys(ctx context.Context, maxResults int32) ([]keyvault.DeletedKeyItem, error)
PurgeDeletedKey(ctx context.Context, keyName string) (bool, error)
RecoverDeletedKey(ctx context.Context, keyName string) (keyvault.KeyBundle, error)*/
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
	return outSign, translateAwsKmsError(err)
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

	return outVerify, translateAwsKmsError(err)
}

func (c *AwsKmsClient) DeleteKey(ctx context.Context, id string) (*kms.DisableKeyOutput, error) {
	outDisable, err := c.client.DisableKey(&kms.DisableKeyInput{
		KeyId: &id,
	})
	return outDisable, translateAwsKmsError(err)
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

func translateAwsKmsError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case kms.ErrCodeAlreadyExistsException:
			return errors.AlreadyExistsError("resource already exists")
		case kms.ErrCodeInternalException:
		case kms.ErrCodeLimitExceededException:
			return errors.InternalError("internal error")
		case kms.ErrCodeIncorrectKeyException:
		case kms.ErrCodeIncorrectKeyMaterialException:
		case kms.ErrCodeInvalidAliasNameException:
		case kms.ErrCodeInvalidCiphertextException:
		case kms.ErrCodeInvalidArnException:
		case kms.ErrCodeInvalidStateException:
			return errors.InvalidParameterError("invalid parameter")
		case kms.ErrCodeNotFoundException:
			return errors.NotFoundError("resource was not found")

		}
	}
	return err
}
