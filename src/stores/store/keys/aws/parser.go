package aws

import (
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"time"
)

const (
	awsKeyID             = "aws-KeyId"
	awsCustomKeyStoreId  = "aws-CustomKeyStoreId"
	awsCloudHsmClusterID = "aws-CloudHsmClusterId"
	awsAccountID         = "aws-AccountId"
	awsARN               = "aws-ARN"
)

func parseAlgorithm(pubKeyInfo *kms.GetPublicKeyOutput) *entities.Algorithm {
	algo := &entities.Algorithm{}
	if pubKeyInfo.KeyUsage != nil && *pubKeyInfo.KeyUsage == kms.KeyUsageTypeSignVerify {
		if *pubKeyInfo.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1 {
			algo.Type = entities.Ecdsa
			algo.EllipticCurve = entities.Secp256k1
		}
	}

	return algo
}

func parseMetadata(describedKey *kms.DescribeKeyOutput) *entities.Metadata {
	// createdAt field always provided
	createdAt := describedKey.KeyMetadata.CreationDate

	deletedAt := &time.Time{}
	if describedKey.KeyMetadata.DeletionDate != nil {
		deletedAt = describedKey.KeyMetadata.DeletionDate
	}

	expireAt := &time.Time{}
	if describedKey.KeyMetadata.ValidTo != nil {
		expireAt = describedKey.KeyMetadata.ValidTo
	}

	return &entities.Metadata{
		Disabled:  !*describedKey.KeyMetadata.Enabled,
		ExpireAt:  *expireAt,
		CreatedAt: *createdAt,
		DeletedAt: *deletedAt,
		UpdatedAt: *createdAt, // Cannot update keys so updatedAt = createdAt
	}
}

func parseAnnotations(keyID string, keyDesc *kms.DescribeKeyOutput) map[string]string {
	annotations := make(map[string]string)

	annotations[awsKeyID] = keyID

	if keyDesc.KeyMetadata.CustomKeyStoreId != nil {
		annotations[awsCustomKeyStoreId] = *keyDesc.KeyMetadata.CustomKeyStoreId
	}
	if keyDesc.KeyMetadata.CloudHsmClusterId != nil {
		annotations[awsCloudHsmClusterID] = *keyDesc.KeyMetadata.CloudHsmClusterId
	}
	if keyDesc.KeyMetadata.AWSAccountId != nil {
		annotations[awsAccountID] = *keyDesc.KeyMetadata.AWSAccountId
	}
	if keyDesc.KeyMetadata.Arn != nil {
		annotations[awsARN] = *keyDesc.KeyMetadata.Arn
	}

	return annotations
}

func toKeyType(alg *entities.Algorithm) (string, error) {
	switch {
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		return kms.CustomerMasterKeySpecEccSecgP256k1, nil
	case alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Bn254:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid key type")
	}
}

func toTags(tags map[string]string) []*kms.Tag {
	var keyTags []*kms.Tag
	for key, value := range tags {
		k, v := key, value
		keyTag := kms.Tag{TagKey: &k, TagValue: &v}
		keyTags = append(keyTags, &keyTag)
	}

	return keyTags
}
