package akv

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
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
	errMsg := "failed to create key"

	kty, err := convertToAKVKeyType(alg)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	crv, err := convertToAKVCurve(alg)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	res, err := s.client.CreateKey(ctx, id, kty, crv, convertToAKVKeyAttr(attr), nil, attr.Tags)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	logger.Info("key was created successfully")
	return parseKeyBundleRes(&res), nil
}

func (s *Store) Import(ctx context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)
	errMsg := "failed to import key"

	iWebKey, err := webImportKey(privKey, alg)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	res, err := s.client.ImportKey(ctx, id, iWebKey, convertToAKVKeyAttr(attr), attr.Tags)
	if err != nil {
		logger.WithError(err).Error(errMsg)
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

func (s *Store) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	logger := s.logger.WithField("id", id)
	expireAt := date.NewUnixTimeFromNanoseconds(expirationDate.UnixNano())
	_, err := s.client.UpdateKey(ctx, id, "", &keyvault.KeyAttributes{
		Expires: &expireAt,
	}, nil, nil)
	if err != nil {
		logger.WithError(err).Error("failed to refresh key")
		return err
	}

	logger.Info("key was refreshed successfully")
	return nil
}

func (s *Store) Delete(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)
	res, err := s.client.DeleteKey(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return nil, err
	}

	logger.Info("key was deleted successfully")
	return parseKeyDeleteBundleRes(&res), nil
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

func (s *Store) Sign(ctx context.Context, id, data string) (string, error) {
	logger := s.logger.WithField("id", id).WithField("data", common.ShortString(data, 5))
	errMsg := "failed to sign data"
	b64Data, err := hexToSha256Base64(data)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return "", err
	}

	kItem, err := s.Get(ctx, id)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return "", err
	}

	algo, err := convertToSignatureAlgo(kItem.Algo)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return "", err
	}

	b64Signature, err := s.client.Sign(ctx, id, "", algo, b64Data)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return "", err
	}

	signature, err := base64ToHex(b64Signature)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return "", errors.InvalidFormatError("expected base64 value. %s", err)
	}

	logger.Debug("data was signed successfully")
	return signature, nil
}

func (s *Store) Encrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.ErrNotImplemented
}

func (s *Store) Decrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.ErrNotImplemented
}
