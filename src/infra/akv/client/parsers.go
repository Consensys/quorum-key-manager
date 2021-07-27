package client

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	PurgeDeletedKeyMethod = "PurgeDeletedKey"
)

func parseErrorResponse(err error) error {
	aerr, ok := err.(autorest.DetailedError)
	if !ok {
		return errors.AKVError("%v", err)
	}

	if rerr, ok := aerr.Original.(*azure.RequestError); ok && rerr.ServiceError.Code == "NotSupported" {
		return errors.NotSupportedError("%v", rerr)
	}

	switch aerr.StatusCode.(int) {
	case http.StatusNotFound:
		return errors.NotFoundError(aerr.Original.Error())
	case http.StatusBadRequest:
		return errors.InvalidFormatError(aerr.Original.Error())
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(aerr.Original.Error())
	case http.StatusConflict:
		if aerr.Method == PurgeDeletedKeyMethod {
			return errors.StatusConflictError(aerr.Original.Error())
		}
		return errors.AlreadyExistsError(aerr.Original.Error())
	default:
		return errors.AKVError(aerr.Original.Error())
	}
}

/**
SECRETS
*/
func parseSecretItem(secretItem *keyvault.SecretItem) *entities.Secret {
	return buildNewSecret(secretItem.ID, nil, secretItem.Tags, secretItem.Attributes)
}

func parseDeletedSecretItem(secretItem *keyvault.DeletedSecretItem) *entities.Secret {
	return buildNewSecret(secretItem.ID, nil, secretItem.Tags, secretItem.Attributes)
}

func parseDeleteSecretBundle(secretBundle *keyvault.DeletedSecretBundle) *entities.Secret {
	return buildNewSecret(secretBundle.ID, secretBundle.Value, secretBundle.Tags, secretBundle.Attributes)
}

func parseSecretBundle(secretBundle *keyvault.SecretBundle) *entities.Secret {
	return buildNewSecret(secretBundle.ID, secretBundle.Value, secretBundle.Tags, secretBundle.Attributes)
}

func buildNewSecret(id, value *string, tags map[string]*string, attributes *keyvault.SecretAttributes) *entities.Secret {
	secret := &entities.Secret{
		Tags:     common.Tomapstr(tags),
		Metadata: &entities.Metadata{},
	}
	if value != nil {
		secret.Value = *value
	}

	if id != nil {
		// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/secrets/<secretName>/<secretVersion>
		// chunks := strings.Split(*id, "/")
		// secret.Metadata.Version = chunks[len(chunks)-1]
		// secret.ID = chunks[len(chunks)-2]
		secret.ID = path.Base(*id)
	}
	if expires := attributes.Expires; expires != nil {
		secret.Metadata.ExpireAt = time.Unix(0, expires.Duration().Nanoseconds()).In(time.UTC)
	}
	if created := attributes.Created; created != nil {
		secret.Metadata.CreatedAt = time.Unix(0, created.Duration().Nanoseconds()).In(time.UTC)
	}
	if updated := attributes.Updated; updated != nil {
		secret.Metadata.UpdatedAt = time.Unix(0, updated.Duration().Nanoseconds()).In(time.UTC)
	}
	if enabled := attributes.Enabled; enabled != nil {
		secret.Metadata.Disabled = !*enabled
	}

	return secret
}

/**
KEYS
*/

func parseKeyBundle(res *keyvault.KeyBundle) *entities.Key {
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

func parseDeletedKeyBundle(res *keyvault.DeletedKeyBundle) *entities.Key {
	return parseKeyBundle(&keyvault.KeyBundle{
		Attributes: res.Attributes,
		Key:        res.Key,
		Tags:       res.Tags,
	})
}

func parseKeyItemBundle(res *keyvault.KeyItem) *entities.Key {
	key := &entities.Key{
		Metadata: &entities.Metadata{
			Disabled:  !*res.Attributes.Enabled,
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Tags: common.Tomapstr(res.Tags),
	}

	key.ID, key.Metadata.Version = parseKeyID(res.Kid)
	return key
}

func parseDeletedKeyItemBundle(res *keyvault.DeletedKeyItem) *entities.Key {
	key := &entities.Key{
		Metadata: &entities.Metadata{
			Disabled:  !*res.Attributes.Enabled,
			CreatedAt: time.Time(*res.Attributes.Created),
			UpdatedAt: time.Time(*res.Attributes.Updated),
		},
		Tags: common.Tomapstr(res.Tags),
	}

	key.ID, key.Metadata.Version = parseKeyID(res.Kid)
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

func convertToAKVOps(ops []entities.CryptoOperation) *[]keyvault.JSONWebKeyOperation {
	akvOps := []keyvault.JSONWebKeyOperation{}

	for _, op := range ops {
		switch op {
		case entities.Encryption:
			akvOps = append(akvOps, keyvault.Encrypt, keyvault.Decrypt)
		case entities.Signing:
			akvOps = append(akvOps, keyvault.Sign, keyvault.Verify)
		}
	}

	return &akvOps
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
