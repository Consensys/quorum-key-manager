package aws

import (
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"time"
)

func algoFromAWSPublicKeyInfo(pubKeyInfo *kms.GetPublicKeyOutput) *entities.Algorithm {
	algo := &entities.Algorithm{}
	if pubKeyInfo == nil {
		return algo
	}
	if pubKeyInfo.KeyUsage != nil && *pubKeyInfo.KeyUsage == kms.KeyUsageTypeSignVerify {

		if *pubKeyInfo.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1 {
			algo.Type = entities.Ecdsa
			algo.EllipticCurve = entities.Secp256k1
		}

	}

	return algo
}

func metadataFromAWSKey(createdKey *kms.CreateKeyOutput) *entities.Metadata {

	createdAt := createdKey.KeyMetadata.CreationDate
	deletedAt := &time.Time{}
	if createdKey.KeyMetadata.DeletionDate != nil {
		deletedAt = createdKey.KeyMetadata.DeletionDate
	}
	expireAt := &time.Time{}
	if createdKey.KeyMetadata.ValidTo != nil {
		expireAt = createdKey.KeyMetadata.ValidTo
	}

	return &entities.Metadata{
		Version:   "",
		Disabled:  !*createdKey.KeyMetadata.Enabled,
		ExpireAt:  *expireAt,
		CreatedAt: *createdAt,
		DeletedAt: *deletedAt,
	}
}

func metadataFromAWSDescribeKey(describedKey *kms.DescribeKeyOutput) *entities.Metadata {
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
		// Nothing equivalent to key version was found
		Version:   "",
		Disabled:  !*describedKey.KeyMetadata.Enabled,
		ExpireAt:  *expireAt,
		CreatedAt: *createdAt,
		DeletedAt: *deletedAt,
	}
}

func fillAwsAnnotations(annotations map[string]string, keyDesc *kms.DescribeKeyOutput) {
	if keyDesc.KeyMetadata.CustomKeyStoreId != nil {
		annotations[awsCustomerKeyStoreID] = *keyDesc.KeyMetadata.CustomKeyStoreId
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
}
