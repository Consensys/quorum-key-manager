package hashicorp

import (
	"context"
	"encoding/base64"
	"path"

	entities2 "github.com/consensys/quorum-key-manager/src/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

const (
	idLabel         = "id"
	curveLabel      = "curve"
	algorithmLabel  = "signing_algorithm"
	tagsLabel       = "tags"
	publicKeyLabel  = "public_key"
	privateKeyLabel = "private_key"
	signatureLabel  = "signature"
	createdAtLabel  = "created_at"
	updatedAtLabel  = "updated_at"
)

type Store struct {
	client hashicorp.PluginClient
	logger log.Logger
}

var _ stores.KeyStore = &Store{}

func New(client hashicorp.PluginClient, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Create(_ context.Context, id string, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	if !s.isSupportedAlgo(alg) {
		errMessage := "invalid or not supported elliptic curve and signing algorithm for Hashicorp key creation"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.NotSupportedError(errMessage)
	}

	res, err := s.client.CreateKey(map[string]interface{}{
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

func (s *Store) Import(_ context.Context, id string, privKey []byte, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	if !s.isSupportedAlgo(alg) {
		errMessage := "invalid or not supported elliptic curve and signing algorithm for Hashicorp key import"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.NotSupportedError(errMessage)
	}

	res, err := s.client.ImportKey(map[string]interface{}{
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

	res, err := s.client.GetKey(id)
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

func (s *Store) List(_ context.Context, _, _ uint64) ([]string, error) {
	res, err := s.client.ListKeys()
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

	var ids []string
	for _, id := range keyIds {
		ids = append(ids, id.(string))
	}

	return ids, nil
}

func (s *Store) Update(_ context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.UpdateKey(id, map[string]interface{}{
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

func (s *Store) ListDeleted(_ context.Context, _, _ uint64) ([]string, error) {
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
	err := s.client.DestroyKey(path.Join(id))
	if err != nil {
		errMessage := "failed to permanently delete Hashicorp key"
		s.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Sign(_ context.Context, id string, data []byte, alg *entities2.Algorithm) ([]byte, error) {
	if !s.isSupportedAlgo(alg) {
		errMessage := "invalid or not supported elliptic curve and signing algorithm for Hashicorp signing"
		s.logger.With("elliptic_curve", alg.EllipticCurve, "signing_algorithm", alg.Type).Error(errMessage)
		return nil, errors.NotSupportedError(errMessage)
	}

	logger := s.logger.With("id", id)

	res, err := s.client.Sign(id, data)
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

func (s *Store) isSupportedAlgo(alg *entities2.Algorithm) bool {
	if alg.Type == entities2.Ecdsa && alg.EllipticCurve == entities2.Secp256k1 {
		return true
	}

	if alg.Type == entities2.Eddsa && alg.EllipticCurve == entities2.Babyjubjub {
		return true
	}

	return false
}
