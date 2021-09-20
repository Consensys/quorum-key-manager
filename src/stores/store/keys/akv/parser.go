package akv

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func convertToAKVOps(ops []entities.CryptoOperation) []keyvault.JSONWebKeyOperation {
	var akvOps []keyvault.JSONWebKeyOperation

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

func convertToAKVKeyAttr(attr *entities.Attributes) *keyvault.KeyAttributes {
	kAttr := &keyvault.KeyAttributes{}
	if attr.TTL.Milliseconds() > 0 {
		ttl := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
		kAttr.Expires = &ttl
	}
	return kAttr
}

func algoFromAKVKeyTypeCrv(kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName) *entities.Algorithm {
	algo := &entities.Algorithm{}
	if kty == keyvault.EC {
		algo.Type = entities.Ecdsa
	}

	if crv == keyvault.P256K {
		algo.EllipticCurve = entities.Secp256k1
	}

	return algo
}

func pubKeyBytes(key *keyvault.JSONWebKey) []byte {
	switch {
	case key.Kty == keyvault.EC && key.Crv == keyvault.P256K:
		xBytes, _ := decodePubKeyBase64(*key.X)
		yBytes, _ := decodePubKeyBase64(*key.Y)
		pKey := ecdsa.PublicKey{X: new(big.Int).SetBytes(xBytes), Y: new(big.Int).SetBytes(yBytes)}
		return crypto.FromECDSAPub(&pKey)
	default:
		return nil
	}

}

func parseKeyBundleRes(res *keyvault.KeyBundle) *entities.Key {
	key := &entities.Key{
		PublicKey: pubKeyBytes(res.Key),
		Algo:      algoFromAKVKeyTypeCrv(res.Key.Kty, res.Key.Crv),
		Metadata: &entities.Metadata{
			Disabled:  !*res.Attributes.Enabled,
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Tags: common.Tomapstr(res.Tags),
	}

	key.ID, key.Metadata.Version = parseKeyID(res.Key.Kid)
	return key
}

func parseKeyID(kid *string) (id, version string) {
	if kid == nil {
		return "", ""
	}
	// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/keys/<keyName>/<secretVersion>
	chunks := strings.Split(*kid, "/")
	var idx int
	for idx = range chunks {
		if chunks[idx] == "keys" {
			break
		}
	}

	if len(chunks) > idx+1 {
		id = chunks[idx+1]
	}
	if len(chunks) > idx+2 {
		version = chunks[idx+2]
	}
	return id, version
}

func parseKeyDeleteBundleRes(res *keyvault.DeletedKeyBundle) *entities.Key {
	return parseKeyBundleRes(&keyvault.KeyBundle{
		Attributes: res.Attributes,
		Key:        res.Key,
		Tags:       res.Tags,
	})
}

func decodePubKeyBase64(src string) ([]byte, error) {
	b := make([]byte, 32)
	for base64.RawURLEncoding.DecodedLen(len(src)) < 32 {
		src += string(base64.StdPadding)
	}

	_, err := base64.RawURLEncoding.Decode(b, []byte(src))
	if err != nil {
		return nil, err
	}
	return b, nil
}
