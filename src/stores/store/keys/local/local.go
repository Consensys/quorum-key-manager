package local

import (
	"context"
	"encoding/base64"

	"github.com/consensys/quorum-key-manager/pkg/crypto/ecdsa"
	"github.com/consensys/quorum-key-manager/pkg/crypto/eddsa"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

type Store struct {
	secretStore stores.SecretStore
	db          database.Secrets
	logger      log.Logger
}

var _ stores.KeyStore = &Store{}

func New(secretStore stores.SecretStore, db database.Secrets, logger log.Logger) *Store {
	return &Store{
		secretStore: secretStore,
		logger:      logger,
		db:          db,
	}
}

func (s *Store) Get(_ context.Context, _ string) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *Store) List(ctx context.Context, _, _ uint64) ([]string, error) {
	var ids []string
	items, err := s.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	return ids, nil
}

func (s *Store) GetDeleted(_ context.Context, _ string) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *Store) ListDeleted(ctx context.Context, _, _ uint64) ([]string, error) {
	var ids []string
	items, err := s.db.GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	return ids, nil
}

func (s *Store) Create(ctx context.Context, id string, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.create(ctx, id, nil, alg, attr)
}

func (s *Store) Import(ctx context.Context, id string, importedPrivKey []byte, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.create(ctx, id, importedPrivKey, alg, attr)
}

func (s *Store) create(ctx context.Context, id string, importedPrivKey []byte, alg *entities2.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id).With("signing_algorithm", alg.Type).With("curve", alg.EllipticCurve)

	var privKey []byte
	var pubKey []byte
	var err error
	switch {
	case alg.Type == entities2.Eddsa && alg.EllipticCurve == entities2.Babyjubjub:
		privKey, pubKey, err = eddsa.CreateBabyjubjub(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate EDDSA/Babyjujub key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	case alg.Type == entities2.Ecdsa && alg.EllipticCurve == entities2.Secp256k1:
		privKey, pubKey, err = ecdsa.CreateSecp256k1(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate Secp256k1/ECDSA key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	case alg.Type == entities2.Eddsa && alg.EllipticCurve == entities2.X25519:
		privKey, pubKey, err = eddsa.CreateX25519(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate EDDSA/Curve25519 key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	secret, err := s.secretStore.Set(ctx, id, base64.StdEncoding.EncodeToString(privKey), attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		secret, err = s.secretStore.Get(ctx, id, "")
	}
	if err != nil {
		return nil, err
	}

	_, err = s.db.Add(ctx, secret)
	if err != nil {
		return nil, err
	}

	return &entities.Key{
		ID:        id,
		PublicKey: pubKey,
		Algo: &entities2.Algorithm{
			Type:          alg.Type,
			EllipticCurve: alg.EllipticCurve,
		},
		Metadata: &entities.Metadata{
			Disabled:  false,
			CreatedAt: secret.Metadata.CreatedAt,
			UpdatedAt: secret.Metadata.UpdatedAt,
		},
		Tags: secret.Tags,
	}, nil
}

func (s *Store) Update(_ context.Context, _ string, _ *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		derr := dbtx.Delete(ctx, id)
		if derr != nil {
			return derr
		}

		derr = s.secretStore.Delete(ctx, id)
		if derr != nil && !errors.IsNotSupportedError(derr) { // If the underlying store does not support deleting, we only delete in DB
			return derr
		}

		return nil
	})
}

func (s *Store) Restore(ctx context.Context, id string) error {
	_, err := s.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	return s.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err := dbtx.Restore(ctx, id)
		if err != nil {
			return err
		}

		err = s.secretStore.Restore(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) {
			return err
		}

		return nil
	})
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	_, err := s.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	return s.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err := dbtx.Purge(ctx, id)
		if err != nil {
			return err
		}

		err = s.secretStore.Destroy(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) {
			return err
		}

		return nil
	})
}

func (s *Store) Sign(ctx context.Context, id string, data []byte, algo *entities2.Algorithm) ([]byte, error) {
	logger := s.logger.With("id", id).With("type", algo.Type).With("curve", algo.EllipticCurve)

	secret, err := s.secretStore.Get(ctx, id, "")
	if err != nil {
		return nil, err
	}

	privkey, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		errMessage := "failed to decode private key secret"
		logger.Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	var signature []byte
	switch {
	case algo.Type == entities2.Eddsa && algo.EllipticCurve == entities2.Babyjubjub:
		signature, err = eddsa.SignBabyjubjub(privkey, data)
	case algo.Type == entities2.Ecdsa && algo.EllipticCurve == entities2.Secp256k1:
		signature, err = ecdsa.SignSecp256k1(privkey, data)
	case algo.Type == entities2.Eddsa && algo.EllipticCurve == entities2.X25519:
		signature, err = eddsa.SignX25519(privkey, data)
	default:
		errMessage := "signing algorithm and curve combination not supported for signing"
		logger.With("algorithm", algo.Type, "curve", algo.EllipticCurve).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	if err != nil {
		s.logger.WithError(err).Error("failed to sign")
	}

	return signature, nil
}

func (s *Store) Verify(_ context.Context, pubKey, data, sig []byte, algo *entities2.Algorithm) error {
	return errors.ErrNotSupported
}

func (s *Store) Encrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
