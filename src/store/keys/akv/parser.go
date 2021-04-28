package akv

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func convertToAKVOps(ops []entities.CryptoOperation) []keyvault.JSONWebKeyOperation {
	akvOps := []keyvault.JSONWebKeyOperation{}

	// TODO Verify list of KM operation and their corresponding into AKV
	for _, op := range ops {
		switch op {
		case entities.Encryption:
			akvOps = append(akvOps, keyvault.Encrypt, keyvault.Decrypt)
		case entities.Signing:
			akvOps = append(akvOps, keyvault.Sign, keyvault.Verify)
		}
	}

	return akvOps
}

func convertToAKVCurve(alg *entities.Algorithm) keyvault.JSONWebKeyCurveName {
	switch alg.EllipticCurve {
	case entities.Secp256k1:
		return keyvault.P256K
	default:
		return ""
	}
}

func convertToAKVKeyType(alg *entities.Algorithm) keyvault.JSONWebKeyType {
	switch alg.Type {
	case entities.Ecdsa:
		return keyvault.EC
	default:
		return ""
	}
}

func convertToSignatureAlgo(alg *entities.Algorithm) keyvault.JSONWebKeySignatureAlgorithm {
	switch alg.Type {
	case entities.Ecdsa:
		switch alg.EllipticCurve {
		case entities.Secp256k1:
			return keyvault.ES256
		default:
			return ""
		}
	default:
		return ""
	}
}

func parseKeyBundleRes(res *keyvault.KeyBundle) *entities.Key {
	key := &entities.Key{
		ID:        *res.Key.Kid,
		PublicKey: *res.Key.X,
		Algo: &entities.Algorithm{
			Type:          string(res.Key.Kty), // @TODO Parse into KM curve type
			EllipticCurve: string(res.Key.Crv), // @TODO Parse into KM curve type
		},
		Metadata: &entities.Metadata{
			Version:   "", // TODO
			Disabled:  !*res.Attributes.Enabled,
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Tags: common.Tomapstr(res.Tags),
	}

	return key
}

func parseKeyDeleteBundleRes(res *keyvault.DeletedKeyBundle) *entities.Key {
	return parseKeyBundleRes(&keyvault.KeyBundle{
		Attributes: res.Attributes,
		Key:        res.Key,
		Tags:       res.Tags,
	})
}
