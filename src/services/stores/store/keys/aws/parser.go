package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/aws/aws-sdk-go/service/kms"
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
