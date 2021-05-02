package akv

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// TODO Verify list of KM operation and their corresponding into AKV
func convertToAKVOps(ops []entities.CryptoOperation) []keyvault.JSONWebKeyOperation {
	akvOps := []keyvault.JSONWebKeyOperation{}

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

func convertToAKVCurve(alg *entities.Algorithm) (keyvault.JSONWebKeyCurveName, error) {
	switch alg.EllipticCurve {
	case entities.Secp256k1:
		return keyvault.P256K, nil
	default:
		return "", errors.NotImplementedError
	}
}

func convertToAKVKeyType(alg *entities.Algorithm) (keyvault.JSONWebKeyType, error) {
	switch alg.Type {
	case entities.Ecdsa:
		return keyvault.EC, nil
	default:
		return "", errors.NotImplementedError
	}
}

func WebImportKey(privKey string, alg *entities.Algorithm) (*keyvault.JSONWebKey, error) {
	var pKeyD, pKeyX, pKeyY string
	switch alg.Type {
	case entities.Ecdsa:
		pKey, err := crypto.HexToECDSA(privKey)
		if err != nil {
			return nil, errors.InvalidFormatError("invalid private key format. %s", err.Error())
		}

		pKeyD = base64.URLEncoding.EncodeToString(pKey.D.Bytes())
		pKeyX = base64.URLEncoding.EncodeToString(pKey.X.Bytes())
		pKeyY = base64.URLEncoding.EncodeToString(pKey.Y.Bytes())
	default:
		return nil, errors.NotImplementedError
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

func algoFromAKVKeyTypeCrv(kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName) *entities.Algorithm {
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

func pubKeyString(key *keyvault.JSONWebKey) string {
	switch {
	case key.Kty == keyvault.EC && key.Crv == keyvault.P256K:
		xBytes := decodeBase64(*key.X, 32)
		yBytes := decodeBase64(*key.Y, 32)
		pKey := ecdsa.PublicKey{X: new(big.Int).SetBytes(xBytes), Y: new(big.Int).SetBytes(yBytes)}
		return hexutil.Encode(crypto.FromECDSAPub(&pKey))
	default:
		return ""
	}

}

func convertToSignatureAlgo(alg *entities.Algorithm) (keyvault.JSONWebKeySignatureAlgorithm, error) {
	switch alg.Type {
	case entities.Ecdsa:
		switch alg.EllipticCurve {
		case entities.Secp256k1:
			return keyvault.ES256K, nil
		default:
			return "", errors.NotImplementedError
		}
	default:
		return "", errors.NotImplementedError
	}
}

func parseKeyBundleRes(res *keyvault.KeyBundle) *entities.Key {
	key := &entities.Key{
		PublicKey: pubKeyString(res.Key),
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

func decodeBase64(src string, n int) []byte {
	b := make([]byte, n)
	for base64.StdEncoding.DecodedLen(len(src)) < n {
		src = src + string(base64.StdPadding)
	}
	base64.URLEncoding.Decode(b, []byte(src))
	return b
}
