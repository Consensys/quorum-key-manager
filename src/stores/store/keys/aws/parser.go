package aws

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"math/big"
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
	key := &entities.Key{
		ID:          id,
		Tags:        tags,
		Annotations: parseAnnotations(*kmsPubKey.KeyId, kmsDescribe),
		Metadata:    &entities.Metadata{},
	}

	switch {
	case *kmsPubKey.KeyUsage == kms.KeyUsageTypeSignVerify && *kmsPubKey.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1:
		key.Algo = &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}

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
	key.Metadata.CreatedAt = *kmsDescribe.KeyMetadata.CreationDate
	if kmsDescribe.KeyMetadata.DeletionDate != nil {
		key.Metadata.DeletedAt = *kmsDescribe.KeyMetadata.DeletionDate
	}
	key.Metadata.Disabled = !*kmsDescribe.KeyMetadata.Enabled

	return key, nil
}

func parseAnnotations(keyID string, keyDesc *kms.DescribeKeyOutput) *entities.Annotation {
	annotations := &entities.Annotation{
		AWSKeyID: keyID,
	}

	if keyDesc.KeyMetadata.CustomKeyStoreId != nil {
		annotations.AWSCustomKeyStoreID = *keyDesc.KeyMetadata.CustomKeyStoreId
	}
	if keyDesc.KeyMetadata.CloudHsmClusterId != nil {
		annotations.AWSCloudHsmClusterID = *keyDesc.KeyMetadata.CloudHsmClusterId
	}
	if keyDesc.KeyMetadata.AWSAccountId != nil {
		annotations.AWSAccountID = *keyDesc.KeyMetadata.AWSAccountId
	}
	if keyDesc.KeyMetadata.Arn != nil {
		annotations.AWSArn = *keyDesc.KeyMetadata.Arn
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
