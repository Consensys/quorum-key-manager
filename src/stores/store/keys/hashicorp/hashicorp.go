package hashicorp

import (
	"context"
	"encoding/base64"
	"path"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

const (
	urlPath         = "keys"
	idLabel         = "id"
	curveLabel      = "curve"
	algorithmLabel  = "signing_algorithm"
	tagsLabel       = "tags"
	publicKeyLabel  = "public_key"
	privateKeyLabel = "private_key"
	dataLabel       = "data"
	signatureLabel  = "signature"
	createdAtLabel  = "created_at"
	updatedAtLabel  = "updated_at"
)

type Store struct {
	client     hashicorp.VaultClient
	mountPoint string
	logger     log.Logger
}

var _ stores.KeyStore = &Store{}

func New(client hashicorp.VaultClient, mountPoint string, logger log.Logger) *Store {
	return &Store{
		client:     client,
		mountPoint: mountPoint,
		logger:     logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(_ context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.Write(s.pathKeys(""), map[string]interface{}{
		idLabel:        id,
		curveLabel:     alg.EllipticCurve,
		algorithmLabel: alg.Type,
		tagsLabel:      attr.Tags,
	})
	if err != nil {
		errMessage := "failed to create Hashicorp key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseAPISecretToKey(res)
}

func (s *Store) Import(_ context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.Write(s.pathKeys("import"), map[string]interface{}{
		idLabel:         id,
		curveLabel:      alg.EllipticCurve,
		algorithmLabel:  alg.Type,
		tagsLabel:       attr.Tags,
		privateKeyLabel: base64.URLEncoding.EncodeToString(privKey),
	})
	if err != nil {
		errMessage := "failed to import Hashicorp key"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseAPISecretToKey(res)
}

func (s *Store) Get(_ context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	res, err := s.client.Read(s.pathKeys(id), nil)
	if err != nil {
		errMessage := "failed to get Hashicorp key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	if res.Data["error"] != nil {
		errMessage := "could not find key pair"
		logger.Error(errMessage)
		return nil, errors.NotFoundError(errMessage)
	}

	return parseAPISecretToKey(res)
}

func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathKeys(""))
	if err != nil {
		errMessage := "failed to list Hashicorp keys"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	if res == nil || res.Data == nil || res.Data["keys"] == nil {
		return []string{}, nil
	}

	keyIds, ok := res.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	var ids = []string{}
	for _, id := range keyIds {
		ids = append(ids, id.(string))
	}

	return ids, nil
}

func (s *Store) Update(_ context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.Write(s.pathKeys(id), map[string]interface{}{
		tagsLabel: attr.Tags,
	})
	if err != nil {
		errMessage := "failed to update Hashicorp key"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseAPISecretToKey(res)
}

func (s *Store) Delete(_ context.Context, _ string) error {
	err := errors.NotSupportedError("delete key is not supported")
	s.logger.Warn(err.Error())
	return err
}

func (s *Store) GetDeleted(_ context.Context, _ string) (*entities.Key, error) {
	err := errors.NotSupportedError("get deleted key is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) ListDeleted(_ context.Context) ([]string, error) {
	err := errors.NotSupportedError("list deleted keys is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) Restore(_ context.Context, _ string) error {
	err := errors.NotSupportedError("restore key is not supported")
	s.logger.Warn(err.Error())
	return err
}

func (s *Store) Destroy(_ context.Context, id string) error {
	err := s.client.Delete(path.Join(s.pathKeys(id), "destroy"), map[string][]string{})
	if err != nil {
		errMessage := "failed to permanently delete Hashicorp key"
		s.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Sign(_ context.Context, id string, data []byte, _ *entities.Algorithm) ([]byte, error) {
	logger := s.logger.With("id", id)

	res, err := s.client.Write(path.Join(s.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: base64.URLEncoding.EncodeToString(data),
	})
	if err != nil {
		errMessage := "failed to sign using Hashicorp key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	signature, err := base64.URLEncoding.DecodeString(res.Data[signatureLabel].(string))
	if err != nil {
		errMessage := "failed to decode signature from Hashicorp Vault"
		logger.WithError(err).Error(errMessage)
		return nil, errors.HashicorpVaultError(errMessage)
	}

	return signature, nil
}

func (s *Store) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	err := errors.NotSupportedError("verify signature is not supported")
	s.logger.Warn(err.Error())
	return err
}

func (s *Store) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) pathKeys(suffix string) string {
	return path.Join(s.mountPoint, urlPath, suffix)
}
