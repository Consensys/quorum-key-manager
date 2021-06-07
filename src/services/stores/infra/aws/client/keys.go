package client

import (
	"context"
	"os"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/aws/aws-sdk-go/service/kms"
)

const (
	aliasPrefix = "alias/"
)

func (c *AwsKmsClient) CreateKey(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*kms.CreateKeyOutput, error) {
	// Always create with same usage for key now (sign & verify)
	keyUsage := kms.KeyUsageTypeSignVerify

	keySpec, err := convertToAWSKeyType(alg)
	if err != nil {
		return nil, err
	}

	outKey, err := c.client.CreateKey(&kms.CreateKeyInput{
		CustomerMasterKeySpec: &keySpec,
		KeyUsage:              &keyUsage,
		Tags:                  extractAWSTags(attr),
	})

	if err != nil {
		return nil, err
	}

	aliasName := aliasPrefix + id
	_, err = c.client.CreateAlias(&kms.CreateAliasInput{
		AliasName:   &aliasName,
		TargetKeyId: outKey.KeyMetadata.KeyId,
	})

	if err != nil {
		return nil, err
	}

	return outKey, nil
}

func (c *AwsKmsClient) GetPublicKey(ctx context.Context, id string) (*kms.GetPublicKeyOutput, error) {
	return c.client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: &id,
	})
}

func (c *AwsKmsClient) ListKeys(ctx context.Context, limit int64, marker string) (*kms.ListKeysOutput, error) {
	input := &kms.ListKeysInput{}
	if limit > 0 {
		input.Limit = &limit
	}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	return c.client.ListKeys(input)
}

func (c *AwsKmsClient) ListTags(ctx context.Context, id, marker string) (*kms.ListResourceTagsOutput, error) {
	input := &kms.ListResourceTagsInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	return c.client.ListResourceTags(input)
}

func (c *AwsKmsClient) ListAliases(ctx context.Context, id, marker string) (*kms.ListAliasesOutput, error) {
	input := &kms.ListAliasesInput{KeyId: &id}
	if len(marker) > 0 {
		input.Marker = &marker
	}
	return c.client.ListAliases(input)
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
	return c.client.Sign(&kms.SignInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		SigningAlgorithm: &signingAlg,
	})
}

func (c *AwsKmsClient) Verify(ctx context.Context, id string, msg, signature []byte) (*kms.VerifyOutput, error) {
	msgType := kms.MessageTypeDigest
	signingAlg := kms.SigningAlgorithmSpecEcdsaSha256
	return c.client.Verify(&kms.VerifyInput{
		KeyId:            &id,
		Message:          msg,
		MessageType:      &msgType,
		Signature:        signature,
		SigningAlgorithm: &signingAlg,
	})
}

func (c *AwsKmsClient) DeleteKey(ctx context.Context, id string) (*kms.DisableKeyOutput, error) {
	return c.client.DisableKey(&kms.DisableKeyInput{
		KeyId: &id,
	})
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

func extractAWSTags(attr *entities.Attributes) []*kms.Tag {
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
