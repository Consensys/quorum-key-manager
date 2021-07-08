package akv

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/log"
	"github.com/consensys/quorum-key-manager/src/stores/infra/akv"
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
	logger := s.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("creating key")

	kty, err := convertToAKVKeyType(alg)
	if err != nil {
		logger.WithError(err).Error("failed to create request: key type")
		return nil, err
	}

	crv, err := convertToAKVCurve(alg)
	if err != nil {
		logger.WithError(err).Error("failed to create request: curve")
		return nil, err
	}

	res, err := s.client.CreateKey(ctx, id, kty, crv, convertToAKVKeyAttr(attr), nil, attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, err
	}

	logger.Info("key was created successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("importing key")

	iWebKey, err := webImportKey(privKey, alg)
	if err != nil {
		logger.WithError(err).Error("failed to create request")
		return nil, err
	}

	res, err := s.client.ImportKey(ctx, id, iWebKey, convertToAKVKeyAttr(attr), attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to import key")
		return nil, err
	}

	logger.Info("key was imported successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	res, err := s.client.GetKey(ctx, id, "")
	if err != nil {
		logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	logger.Debug("key retrieved successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	res, err := s.client.GetKeys(ctx, 0)
	if err != nil {
		s.logger.WithError(err).Error("failed to list keys")
		return nil, err
	}

	kIDs := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIDs = append(kIDs, kID)
	}

	s.logger.Debug("keys listed successfully")
	return kIDs, nil
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)
	logger.Debug("updating key")

	expireAt := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
	res, err := s.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, convertToAKVOps(attr.Operations), attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to update key")
		return nil, err
	}

	logger.Info("key updated successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("deleting key")

	_, err := s.client.DeleteKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	res, err := s.client.GetDeletedKey(ctx, id)
	if err != nil {
		logger.Error("failed to get deleted key")
		return nil, err
	}

	logger.Debug("deleted key retrieved successfully")
	return parseKeyDeleteBundleRes(&res), nil
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	res, err := s.client.GetDeletedKeys(ctx, 0)
	if err != nil {
		s.logger.Error("failed to list deleted keys")
		return nil, err
	}

	kIds := []string{}
	for _, kItem := range res {
		kID, _ := parseKeyID(kItem.Kid)
		kIds = append(kIds, kID)
	}

	s.logger.Debug("deleted keys listed successfully")
	return kIds, nil
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("restoring key")

	_, err := s.client.RecoverDeletedKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to restore key")
		return err
	}

	logger.Info("key restored successfully")
	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("destroying key")

	_, err := s.client.PurgeDeletedKey(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("failed to permanently delete key")
		return err
	}

	logger.Info("key was permanently deleted")
	return nil
}

func (s *Store) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := s.logger.With("id", id)
	logger.Debug("signing payload")

	kItem, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	algo, err := convertToSignatureAlgo(kItem.Algo)
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return nil, err
	}

	b64Signature, err := s.client.Sign(ctx, id, "", algo, base64.StdEncoding.EncodeToString(data))
	if err != nil {
		errMessage := "failed to sign payload"
		logger.WithError(err).Error(errMessage)
		return nil, err
	}

	signature, err := base64.RawURLEncoding.DecodeString(b64Signature)
	if err != nil {
		errMessage := "failed to decode signature from AKV vault"
		logger.WithError(err).Error(errMessage)
		return nil, errors.AKVError(errMessage)
	}

	logger.Debug("payload signed successfully")
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
