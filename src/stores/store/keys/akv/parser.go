package akv

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"
	"strings"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/common"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func convertToAKVKeyAttr(attr *entities.Attributes) *keyvault.KeyAttributes {
	kAttr := &keyvault.KeyAttributes{}
	if attr.TTL.Milliseconds() > 0 {
		ttl := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
		kAttr.Expires = &ttl
	}
	return kAttr
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
		ID:        parseKeyID(res.Key.Kid),
		PublicKey: pubKeyBytes(res.Key),
		Tags:      common.Tomapstr(res.Tags),
		Metadata: &entities.Metadata{
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Algo: &entities.Algorithm{},
	}

	if res.Key.Kty == keyvault.EC {
		key.Algo.Type = entities.Ecdsa
	}

	if res.Key.Crv == keyvault.P256K {
		key.Algo.EllipticCurve = entities.Secp256k1
	}

	return key
}

func parseKeyID(kid *string) string {
	// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/keys/<keyName>/<secretVersion>
	chunks := strings.Split(*kid, "/")
	var idx int
	for idx = range chunks {
		if chunks[idx] == "keys" {
			break
		}
	}

	if len(chunks) > idx+1 {
		return chunks[idx+1]
	}

	return ""
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
