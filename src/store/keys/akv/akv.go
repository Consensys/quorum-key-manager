package akv

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
)

type Store struct {
	client akv.KeysClient
	logger *log.Logger
}

var _ keys.Store = &Store{}

func New(client akv.KeysClient, logger *log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)

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
	logger := s.logger.WithField("id", id)

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
	logger := s.logger.WithField("id", id)
	res, err := s.client.GetKey(ctx, id, "")
	if err != nil {
		logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	logger.Debug("key was retrieved successfully")
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

	s.logger.Debug("keys were listed successfully")
	return kIDs, nil
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)
	expireAt := date.NewUnixTimeFromNanoseconds(time.Now().Add(attr.TTL).UnixNano())
	res, err := s.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, convertToAKVOps(attr.Operations), attr.Tags)
	if err != nil {
		logger.WithError(err).Error("failed to update key")
		return nil, err
	}

	logger.Info("key was updated successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)
	_, err := s.client.DeleteKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return err
	}

	logger.Info("key was deleted successfully")
	return nil
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	res, err := s.client.GetDeletedKey(ctx, id)
	if err != nil {
		s.logger.WithField("id", id).Error("failed to get deleted keys")
		return nil, err
	}

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

	return kIds, nil
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)
	_, err := s.client.RecoverDeletedKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to undelete key")
		return err
	}

	logger.Info("key was undeleted successfully")
	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)
	_, err := s.client.PurgeDeletedKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to destroy key")
		return err
	}

	logger.Info("key was destroyed successfully")
	return nil
}

func (s *Store) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := s.logger.WithField("id", id)

	kItem, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	algo, err := convertToSignatureAlgo(kItem.Algo)
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return nil, err
	}

	b64Signature, err := s.client.Sign(ctx, id, "", algo, base64.RawURLEncoding.EncodeToString(crypto.Keccak256(data)))
	if err != nil {
		errMessage := "failed to sign payload"
		logger.WithError(err).Error(errMessage)
		return nil, err
	}

	signature, err := base64.RawURLEncoding.DecodeString(b64Signature)
	if err != nil {
		errMessage := "failed to decode signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.AKVConnectionError(errMessage)
	}

	logger.Debug("data was signed successfully")
	return signature, nil
}

func (s *Store) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
