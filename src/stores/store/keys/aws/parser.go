package aws

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

const (
	awsKeyID             = "aws-KeyId"
	awsCustomKeyStoreID  = "aws-CustomKeyStoreId"
	awsCloudHsmClusterID = "aws-CloudHsmClusterId"
	awsAccountID         = "aws-AccountId"
	awsARN               = "aws-ARN"
)

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type signatureInfo struct {
	R, S *big.Int
}

func parseKey(id string, kmsPubKey *kms.GetPublicKeyOutput, kmsDescribe *kms.DescribeKeyOutput, tags map[string]string) (*entities.Key, error) {
	var algo *entities.Algorithm
	var pubKey []byte

	switch {
	case *kmsPubKey.KeyUsage == kms.KeyUsageTypeSignVerify && *kmsPubKey.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1:
		algo = &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}

		val := &publicKeyInfo{}
		_, err := asn1.Unmarshal(kmsPubKey.PublicKey, val)
		if err != nil {
			return nil, err
		}
		pubKey = val.PublicKey.Bytes
	default:
		return nil, fmt.Errorf("unsupported public key type returned from AWS KMS")
	}

	return &entities.Key{
		ID:          id,
		PublicKey:   pubKey,
		Algo:        algo,
		Metadata:    parseMetadata(kmsDescribe),
		Tags:        tags,
		Annotations: parseAnnotations(*kmsPubKey.KeyId, kmsDescribe),
	}, nil
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
		Version:   "1",
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
		annotations[awsCustomKeyStoreID] = *keyDesc.KeyMetadata.CustomKeyStoreId
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

func parseSignature(kmsSign *kms.SignOutput) ([]byte, error) {
	val := &signatureInfo{}
	_, err := asn1.Unmarshal(kmsSign.Signature, val)
	if err != nil {
		return nil, err
	}

	// ensure signature size is 64
	sig := make([]byte, 64)
	// copy R in first half
	copy(sig[len(sig)/2-len(val.R.Bytes()):], val.R.Bytes()[:])
	// copy S in second half
	copy(sig[len(sig)-len(val.S.Bytes()):], val.S.Bytes()[:])

	return sig, nil
}

func toTags(tags map[string]string) []*kms.Tag {
	var keyTags []*kms.Tag
	for key, value := range tags {
		k, v := key, value
		keyTags = append(keyTags, &kms.Tag{TagKey: &k, TagValue: &v})
	}

	return keyTags
}
