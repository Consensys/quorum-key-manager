package hashicorp

import (
	"context"
	"encoding/base64"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"path"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

const (
	urlPath         = "keys"
	idLabel         = "id"
	curveLabel      = "curve"
	algorithmLabel  = "algorithm"
	tagsLabel       = "tags"
	publicKeyLabel  = "public_key"
	privateKeyLabel = "privateKey"
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

var _ keys.Store = &Store{}

func New(client hashicorp.VaultClient, mountPoint string, logger log.Logger) *Store {
	return &Store{
		client:     client,
		mountPoint: mountPoint,
		logger:     logger,
	}
}

func (s *Store) Create(_ context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*models.Key, error) {
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

	return parseResponse(res)
}

func (s *Store) Import(_ context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*models.Key, error) {
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

	return parseResponse(res)
}

func (s *Store) Update(_ context.Context, id string, attr *entities.Attributes) (*models.Key, error) {
	res, err := s.client.Write(s.pathKeys(id), map[string]interface{}{
		tagsLabel: attr.Tags,
	})
	if err != nil {
		errMessage := "failed to update Hashicorp key"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return parseResponse(res)
}

func (s *Store) Delete(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *Store) Undelete(_ context.Context, _ string) error {
	return errors.ErrNotSupported
}

func (s *Store) Destroy(_ context.Context, id string) error {
	err := s.client.Delete(path.Join(s.pathKeys(id), "destroy"))
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

func (s *Store) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) pathKeys(suffix string) string {
	return path.Join(s.mountPoint, urlPath, suffix)
}
