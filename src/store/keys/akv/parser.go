package akv

import (
	"strings"
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

func buildAlgoFromAKVKeyTypeCrv(kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName) *entities.Algorithm  {
	algo := &entities.Algorithm{}
	switch kty {
	case keyvault.EC:
		algo.Type = entities.Ecdsa
	}
	
	switch crv {
	case keyvault.P256K:
		algo.EllipticCurve = entities.Secp256k1
	}
	
	return algo
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
		PublicKey: *res.Key.X,
		Algo: buildAlgoFromAKVKeyTypeCrv(res.Key.Kty, res.Key.Crv),
		Metadata: &entities.Metadata{
			Disabled:  !*res.Attributes.Enabled,
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Tags: common.Tomapstr(res.Tags),
	}

	if res.Key.Kid != nil {
		// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/keys/<keyName>/<secretVersion>
		chunks := strings.Split(*res.Key.Kid, "/")
		key.Metadata.Version = chunks[len(chunks)-1]
		key.ID = chunks[len(chunks)-2]
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
