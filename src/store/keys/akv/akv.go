package akv

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	akvclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
)

// Store is an implementation of key store relying on Hashicorp Vault ConsenSys secret engine
type KeyStore struct {
	client akv.KeysClient
}

var _ keys.Store = KeyStore{}

// New creates an HashiCorp key store
func New(client akv.KeysClient) *KeyStore {
	return &KeyStore{
		client: client,
	}
}

func (k KeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

func (k KeyStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	expireAt := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
	res, err := k.client.CreateKey(ctx, id, convertToAKVKeyType(alg), convertToAKVCurve(alg), &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, convertToAKVOps(attr.Operations), attr.Tags)

	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyBundleRes(&res), nil
}

func (k KeyStore) Import(ctx context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	kOps := []string{}
	for _, op := range convertToAKVOps(attr.Operations) {
		kOps = append(kOps, string(op))
	}
	res, err := k.client.ImportKey(ctx, id, &keyvault.JSONWebKey{
		Crv: convertToAKVCurve(alg),
		Kty: convertToAKVKeyType(alg),
		KeyOps: &kOps,
		K: &privKey,
	}, attr.Tags)

	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyBundleRes(&res), nil
}

func (k KeyStore) Get(ctx context.Context, id, version string) (*entities.Key, error) {
	res, err := k.client.GetKey(ctx, id, version)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyBundleRes(&res), nil
}

func (k KeyStore) List(ctx context.Context) ([]string, error) {
	res, err := k.client.GetKeys(ctx, 0)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	kIds := []string{}
	for _, kItem := range res {
		kIds = append(kIds, *kItem.Kid)
	}
	return kIds, nil
}

func (k KeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	expireAt := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
	// @TODO CHeck if empty version updates latest key
	res, err := k.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, convertToAKVOps(attr.Operations), attr.Tags)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyBundleRes(&res), nil
}

func (k KeyStore) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	expireAt := date.NewUnixTimeFromNanoseconds(expirationDate.UnixNano())
	// @TODO CHeck if empty version updates latest key
	_, err := k.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	},nil, nil)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}
	
	return nil
}

func (k KeyStore) Delete(ctx context.Context, id string) (*entities.Key, error) {
	res, err := k.client.DeleteKey(ctx, id)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyDeleteBundleRes(&res), nil
}

func (k KeyStore) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	res, err := k.client.GetDeletedKey(ctx, id)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}
	
	return parseKeyDeleteBundleRes(&res), nil
}

func (k KeyStore) ListDeleted(ctx context.Context) ([]string, error) {
	res, err := k.client.GetDeletedKeys(ctx, 0)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	kIds := []string{}
	for _, kItem := range res {
		kIds = append(kIds, *kItem.Kid)
	}

	return kIds, nil
}

func (k KeyStore) Undelete(ctx context.Context, id string) error {
	_, err := k.client.RecoverDeletedKey(ctx, id)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}

func (k KeyStore) Destroy(ctx context.Context, id string) error {
	_, err := k.client.PurgeDeletedKey(ctx, id)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}

func (k KeyStore) Sign(ctx context.Context, id, data, version string) (string, error) {
	kItem, err := k.Get(ctx, id, version)
	if err != nil {
		return "", err
	}

	signature, err := k.client.Sign(ctx, id, version, convertToSignatureAlgo(kItem.Algo), data)
	if err != nil {
		return "", akvclient.ParseErrorResponse(err)
	}
	return signature, nil
}

func (k KeyStore) Encrypt(ctx context.Context, id, version, data string) (string, error) {
	//@TODO Shouldn't we expect a encryption algo as argument?
	signature, err := k.client.Encrypt(ctx, id, version, keyvault.RSA15, data)
	if err != nil {
		return "", akvclient.ParseErrorResponse(err)
	}
	return signature, nil
}

func (k KeyStore) Decrypt(ctx context.Context, id, version, data string) (string, error) {
	//@TODO Shouldn't we expect a encryption algo as argument?
	decryptedData, err := k.client.Decrypt(ctx, id, version, keyvault.RSA15, data)
	if err != nil {
		return "", akvclient.ParseErrorResponse(err)
	}

	return decryptedData, nil
}
