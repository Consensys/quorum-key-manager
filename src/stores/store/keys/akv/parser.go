package akv

import (
	"crypto/ecdsa"
	"encoding/base64"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"math/big"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func convertToAKVOps(ops []entities2.CryptoOperation) []keyvault.JSONWebKeyOperation {
	akvOps := []keyvault.JSONWebKeyOperation{}

	for _, op := range ops {
		switch op {
		case entities2.Encryption:
			akvOps = append(akvOps, keyvault.Encrypt, keyvault.Decrypt)
		case entities2.Signing:
			akvOps = append(akvOps, keyvault.Sign, keyvault.Verify)
		}
	}

	return akvOps
}

func convertToAKVCurve(alg *entities2.Algorithm) (keyvault.JSONWebKeyCurveName, error) {
	switch alg.EllipticCurve {
	case entities2.Secp256k1:
		return keyvault.P256K, nil
	case entities2.Bn254:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid elliptic curve")
	}
}

func convertToAKVKeyType(alg *entities2.Algorithm) (keyvault.JSONWebKeyType, error) {
	switch alg.Type {
	case entities2.Ecdsa:
		return keyvault.EC, nil
	case entities2.Eddsa:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid key type")
	}
}

func convertToAKVKeyAttr(attr *entities2.Attributes) *keyvault.KeyAttributes {
	kAttr := &keyvault.KeyAttributes{}
	if attr.TTL.Milliseconds() > 0 {
		ttl := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
		kAttr.Expires = &ttl
	}
	return kAttr
}

func webImportKey(privKey []byte, alg *entities2.Algorithm) (*keyvault.JSONWebKey, error) {
	var pKeyD, pKeyX, pKeyY string
	switch alg.Type {
	case entities2.Ecdsa:
		pKey, err := crypto.ToECDSA(privKey)
		if err != nil {
			return nil, errors.InvalidParameterError("invalid private key. %s", err.Error())
		}

		pKeyD = base64.RawURLEncoding.EncodeToString(pKey.D.Bytes())
		pKeyX = base64.RawURLEncoding.EncodeToString(pKey.X.Bytes())
		pKeyY = base64.RawURLEncoding.EncodeToString(pKey.Y.Bytes())
	case entities2.Eddsa:
		return nil, errors.ErrNotSupported
	default:
		return nil, errors.InvalidParameterError("invalid key type")
	}

	var err error
	var crv keyvault.JSONWebKeyCurveName
	var kty keyvault.JSONWebKeyType
	if crv, err = convertToAKVCurve(alg); err != nil {
		return nil, err
	}
	if kty, err = convertToAKVKeyType(alg); err != nil {
		return nil, err
	}

	return &keyvault.JSONWebKey{
		Crv: crv,
		Kty: kty,
		D:   &pKeyD,
		X:   &pKeyX,
		Y:   &pKeyY,
	}, nil
}

func algoFromAKVKeyTypeCrv(kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName) *entities2.Algorithm {
	algo := &entities2.Algorithm{}
	if kty == keyvault.EC {
		algo.Type = entities2.Ecdsa
	}

	if crv == keyvault.P256K {
		algo.EllipticCurve = entities2.Secp256k1
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

func convertToSignatureAlgo(alg *entities2.Algorithm) (keyvault.JSONWebKeySignatureAlgorithm, error) {
	switch alg.Type {
	case entities2.Ecdsa:
		switch alg.EllipticCurve {
		case entities2.Secp256k1:
			return keyvault.ES256K, nil
		case entities2.Bn254:
			return "", errors.ErrNotSupported
		default:
			return "", errors.InvalidParameterError("invalid elliptic curve")
		}
	case entities2.Eddsa:
		return "", errors.ErrNotSupported
	default:
		return "", errors.InvalidParameterError("invalid key type")
	}
}

func parseKeyBundleRes(res *keyvault.KeyBundle) *entities2.Key {
	key := &entities2.Key{
		PublicKey: pubKeyBytes(res.Key),
		Algo:      algoFromAKVKeyTypeCrv(res.Key.Kty, res.Key.Crv),
		Metadata: &entities2.Metadata{
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

func parseKeyDeleteBundleRes(res *keyvault.DeletedKeyBundle) *entities2.Key {
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
