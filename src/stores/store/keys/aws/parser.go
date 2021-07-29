package aws

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"math/big"
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

func parseKey(id string, kmsPubKey *kms.GetPublicKeyOutput, kmsDescribe *kms.DescribeKeyOutput, tags map[string]string) (*models.Key, error) {
	key := &models.Key{
		ID:          id,
		Tags:        tags,
		Annotations: parseAnnotations(*kmsPubKey.KeyId, kmsDescribe),
	}

	switch {
	case *kmsPubKey.KeyUsage == kms.KeyUsageTypeSignVerify && *kmsPubKey.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1:
		key.SigningAlgorithm = string(entities.Ecdsa)
		key.EllipticCurve = string(entities.Secp256k1)

		val := &publicKeyInfo{}
		_, err := asn1.Unmarshal(kmsPubKey.PublicKey, val)
		if err != nil {
			return nil, err
		}
		key.PublicKey = val.PublicKey.Bytes
	default:
		return nil, fmt.Errorf("unsupported public key type returned from AWS KMS")
	}

	// createdAt field always provided
	key.CreatedAt = *kmsDescribe.KeyMetadata.CreationDate
	if kmsDescribe.KeyMetadata.DeletionDate != nil {
		key.DeletedAt = *kmsDescribe.KeyMetadata.DeletionDate
	}
	key.Disabled = !*kmsDescribe.KeyMetadata.Enabled

	return key, nil
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

	return append(val.R.Bytes(), val.S.Bytes()...), nil
}

func toTags(tags map[string]string) []*kms.Tag {
	var keyTags []*kms.Tag
	for key, value := range tags {
		k, v := key, value
		keyTags = append(keyTags, &kms.Tag{TagKey: &k, TagValue: &v})
	}

	return keyTags
}
