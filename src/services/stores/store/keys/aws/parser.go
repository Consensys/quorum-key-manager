package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/aws/aws-sdk-go/service/kms"
)

func algoFromAWSPublicKeyInfo(pubKeyInfo *kms.GetPublicKeyOutput) *entities.Algorithm {
	algo := &entities.Algorithm{}
	if pubKeyInfo == nil {
		return algo
	}
	if pubKeyInfo.KeyUsage != nil && *pubKeyInfo.KeyUsage == kms.KeyUsageTypeSignVerify {

		if *pubKeyInfo.CustomerMasterKeySpec == kms.CustomerMasterKeySpecEccSecgP256k1 {
			algo.Type = kms.CustomerMasterKeySpecEccSecgP256k1
			algo.Type = kms.CustomerMasterKeySpecEccSecgP256k1
		}

	}

	return algo
}
