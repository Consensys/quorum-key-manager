package akv

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

type Store struct {
	client akv.KeysClient
	logger log.Logger
}

var _ keys.Store = &Store{}

func New(client akv.KeysClient, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	var kty keyvault.JSONWebKeyType
	var crv keyvault.JSONWebKeyCurveName

	switch {
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		kty = keyvault.EC
		crv = keyvault.P256K
	default:
		errMessage := "invalid or not supported elliptic curve and signing algorithm for AKV key creation"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	res, err := s.client.CreateKey(ctx, id, kty, crv, convertToAKVKeyAttr(attr), nil, attr.Tags)
	if err != nil {
		errMessage := "failed to create AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseKeyBundleRes(&res), nil
}

func (s *Store) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	var pKeyD, pKeyX, pKeyY string
	var kty keyvault.JSONWebKeyType
	var crv keyvault.JSONWebKeyCurveName

	switch {
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		pKey, err := crypto.ToECDSA(privKey)
		if err != nil {
			errMessage := "invalid private key"
			s.logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		pKeyD = base64.RawURLEncoding.EncodeToString(pKey.D.Bytes())
		pKeyX = base64.RawURLEncoding.EncodeToString(pKey.X.Bytes())
		pKeyY = base64.RawURLEncoding.EncodeToString(pKey.Y.Bytes())
		kty = keyvault.EC
		crv = keyvault.P256K
	default:
		errMessage := "invalid signing algorithm and curve combination for import"
		s.logger.With("signing_algorithm", alg.Type, "elliptic_curve", alg.EllipticCurve).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	iWebKey := &keyvault.JSONWebKey{
		Crv: crv,
		Kty: kty,
		D:   &pKeyD,
		X:   &pKeyX,
		Y:   &pKeyY,
	}
	res, err := s.client.ImportKey(ctx, id, iWebKey, convertToAKVKeyAttr(attr), attr.Tags)
	if err != nil {
		errMessage := "failed to import AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseKeyBundleRes(&res), nil
}

func (s *Store) Get(ctx context.Context, id string) (*entities.Key, error) {
	res, err := s.client.GetKey(ctx, id, "")
	if err != nil {
		errMessage := "failed to get AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseKeyBundleRes(&res), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	res, err := s.client.GetKeys(ctx, 0)
	if err != nil {
		errMessage := "failed to list AKV keys"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	kIDs := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIDs = append(kIDs, kID)
	}

	return kIDs, nil
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	expireAt := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
	res, err := s.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, convertToAKVOps(attr.Operations), attr.Tags)
	if err != nil {
		errMessage := "failed to update AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseKeyBundleRes(&res), nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.client.DeleteKey(ctx, id)
	if err != nil {
		errMessage := "failed to delete AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	res, err := s.client.GetDeletedKey(ctx, id)
	if err != nil {
		errMessage := "failed to get deleted AKV key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseKeyDeleteBundleRes(&res), nil
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	res, err := s.client.GetDeletedKeys(ctx, 0)
	if err != nil {
		errMessage := "failed to list deleted AKV keys"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	kIds := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIds = append(kIds, kID)
	}

	return kIds, nil
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	_, err := s.client.RecoverDeletedKey(ctx, id)
	if err != nil {
		errMessage := "failed to restore AKV key"
		s.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	_, err := s.client.PurgeDeletedKey(ctx, id)
	if err != nil {
		errMessage := "failed to permantently delete AKV key"
		s.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	kItem, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	var algo keyvault.JSONWebKeySignatureAlgorithm
	switch {
	case kItem.Algo.EllipticCurve == entities.Secp256k1 && kItem.Algo.Type == entities.Ecdsa:
		algo = keyvault.ES256K
	default:
		errMessage := "invalid elliptic curve and signing algorithm combination for signing"
		s.logger.With("id", id, "signing_algorithm", kItem.Algo.Type, "elliptic_curve", kItem.Algo.EllipticCurve).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	b64Signature, err := s.client.Sign(ctx, id, "", algo, base64.StdEncoding.EncodeToString(data))
	if err != nil {
		errMessage := "failed to sign using AKV key"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	signature, err := base64.RawURLEncoding.DecodeString(b64Signature)
	if err != nil {
		errMessage := "failed to decode signature from AKV vault"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.AKVError(errMessage)
	}

	return signature, nil
}

func (s *Store) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	return keys.VerifySignature(s.logger, pubKey, data, sig, algo)
}

func (s *Store) Encrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
