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

type Store struct {
	client akv.KeysClient
}

var _ keys.Store = Store{}

func New(client akv.KeysClient) *Store {
	return &Store{
		client: client,
	}
}

func (k Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (k Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {

	kty, err := convertToAKVKeyType(alg)
	if err != nil {
		return nil, err
	}

	crv, err := convertToAKVCurve(alg)
	if err != nil {
		return nil, err
	}

	res, err := k.client.CreateKey(ctx, id, kty, crv, convertToAKVKeyAttr(attr), nil, attr.Tags)

	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseKeyBundleRes(&res), nil
}

func (k Store) Import(ctx context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	iWebKey, err := webImportKey(privKey, alg)
	if err != nil {
		return nil, err
	}

	res, err := k.client.ImportKey(ctx, id, iWebKey, convertToAKVKeyAttr(attr), attr.Tags)

	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseKeyBundleRes(&res), nil
}

func (k Store) Get(ctx context.Context, id, version string) (*entities.Key, error) {
	res, err := k.client.GetKey(ctx, id, version)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseKeyBundleRes(&res), nil
}

func (k Store) List(ctx context.Context) ([]string, error) {
	res, err := k.client.GetKeys(ctx, 0)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	kIDs := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIDs = append(kIDs, kID)
	}
	return kIDs, nil
}

func (k Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
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

func (k Store) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	expireAt := date.NewUnixTimeFromNanoseconds(expirationDate.UnixNano())
	// @TODO CHeck if empty version updates latest key
	_, err := k.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, nil, nil)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}

func (k Store) Delete(ctx context.Context, id string) (*entities.Key, error) {
	res, err := k.client.DeleteKey(ctx, id)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseKeyDeleteBundleRes(&res), nil
}

func (k Store) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	res, err := k.client.GetDeletedKey(ctx, id)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	return parseKeyDeleteBundleRes(&res), nil
}

func (k Store) ListDeleted(ctx context.Context) ([]string, error) {
	res, err := k.client.GetDeletedKeys(ctx, 0)
	if err != nil {
		return nil, akvclient.ParseErrorResponse(err)
	}

	kIds := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIds = append(kIds, kID)
	}

	return kIds, nil
}

func (k Store) Undelete(ctx context.Context, id string) error {
	_, err := k.client.RecoverDeletedKey(ctx, id)
	if err != nil {
		return akvclient.ParseErrorResponse(err)
	}

	return nil
}

func (k Store) Destroy(ctx context.Context, id string) error {
	_, err := k.client.PurgeDeletedKey(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (k Store) Sign(ctx context.Context, id, data, version string) (string, error) {
	b64Data, err := hexToSha256Base64(data)
	if err != nil {
		return "", err
	}

	kItem, err := k.Get(ctx, id, version)
	if err != nil {
		return "", err
	}

	algo, err := convertToSignatureAlgo(kItem.Algo)
	if err != nil {
		return "", err
	}

	b64Signature, err := k.client.Sign(ctx, id, version, algo, b64Data)
	if err != nil {
		return "", akvclient.ParseErrorResponse(err)
	}

	signature, err := base64ToHex(b64Signature)
	if err != nil {
		return "", errors.InvalidFormatError("expected base64 value. %s", err)
	}

	return signature, nil
}

func (k Store) Encrypt(ctx context.Context, id, version, data string) (string, error) {
	return "", errors.ErrNotImplemented
}

func (k Store) Decrypt(ctx context.Context, id, version, data string) (string, error) {
	return "", errors.ErrNotImplemented
}
