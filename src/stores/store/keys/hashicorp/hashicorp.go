package hashicorp

import (
	"context"
	"encoding/base64"
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
	versionLabel    = "version"
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
		return nil, err
	}

	return parseResponse(res)
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
		return nil, err
	}

	return parseResponse(res)
}

func (s *Store) Get(_ context.Context, id string) (*entities.Key, error) {
	res, err := s.client.Read(s.pathKeys(id), nil)
	if err != nil {
		return nil, err
	}

	if res.Data["error"] != nil {
		return nil, errors.NotFoundError("could not find key pair")
	}

	return parseResponse(res)
}

func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathKeys(""))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return parseResponse(res)
}

func (s *Store) Delete(_ context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("destroying key")

	err := s.client.Delete(path.Join(s.pathKeys(id), "destroy"))
	if err != nil {
		s.logger.WithError(err).Error("failed to permanently delete key")
		return err
	}

	logger.Info("key permanently deleted")
	return nil
}

func (s *Store) Sign(_ context.Context, id string, data []byte) ([]byte, error) {
	logger := s.logger.With("id", id)
	logger.Debug("signing payload")

	res, err := s.client.Write(path.Join(s.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: base64.URLEncoding.EncodeToString(data),
	})
	if err != nil {
		logger.WithError(err).Error("failed to sign payload")
		return nil, err
	}

	signature, err := base64.URLEncoding.DecodeString(res.Data[signatureLabel].(string))
	if err != nil {
		errMessage := "failed to decode signature from Hashicorp Vault"
		logger.WithError(err).Error(errMessage)
		return nil, errors.HashicorpVaultError(errMessage)
	}

	logger.Debug("payload signed successfully")
	return signature, nil
}

func (s *Store) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	return keys.VerifySignature(s.logger, pubKey, data, sig, algo)
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
