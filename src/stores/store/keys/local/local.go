package local

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/store/database/models"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/gnark-crypto/crypto/hash"
	eddsabn254 "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

type Store struct {
	secretStore secrets.Store
	db          database.Database
	logger      log.Logger
}

var _ keys.Store = &Store{}

func New(secretStore secrets.Store, db database.Database, logger log.Logger) *Store {
	return &Store{
		secretStore: secretStore,
		db:          db,
		logger:      logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.createKey(ctx, id, nil, alg, attr)
}

func (s *Store) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.createKey(ctx, id, privKey, alg, attr)
}

func (s *Store) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	key, err := s.db.Keys().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug("key retrieved successfully")
	return key.ToEntity(), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	ids := []string{}
	keysRetrieved, err := s.db.Keys().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keysRetrieved {
		ids = append(ids, key.ID)
	}

	s.logger.Debug("keys listed successfully")
	return ids, nil
}

func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id)
	logger.Debug("updating key")

	key, err := s.db.Keys().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	key.Tags = attr.Tags

	err = s.db.Keys().Update(ctx, key)
	if err != nil {
		return nil, err
	}

	logger.Info("key updated successfully")
	return key.ToEntity(), nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("deleting key")

	err := s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		err := s.db.Keys().Remove(ctx, id)
		if err != nil {
			return err
		}

		return s.secretStore.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}

func (s *Store) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	logger := s.logger.With("id", id)

	key, err := s.db.Keys().GetDeleted(ctx, id)
	if err != nil {
		logger.Error("failed to get deleted key")
		return nil, err
	}

	logger.Debug("deleted key retrieved successfully")
	return key.ToEntity(), nil
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	ids := []string{}
	keysRetrieved, err := s.db.Keys().GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keysRetrieved {
		ids = append(ids, key.ID)
	}

	s.logger.Debug("deleted keys listed successfully")
	return ids, nil
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("restoring key")

	key, err := s.db.Keys().GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := s.db.Keys().Restore(ctx, key)
		if derr != nil {
			return derr
		}

		return s.secretStore.Undelete(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key restored successfully")
	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("destroying key")

	_, err := s.db.Keys().GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := s.db.Keys().Purge(ctx, id)
		if derr != nil {
			return derr
		}

		return s.secretStore.Destroy(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key was permanently deleted")
	return nil
}

func (s *Store) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := s.logger.With("id", id)
	logger.Debug("signing payload")

	key, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

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

	switch {
	case key.Algo.Type == entities.Eddsa && key.Algo.EllipticCurve == entities.Bn254:
		return s.signEDDSA(privkey, data)
	case key.Algo.Type == entities.Ecdsa && key.Algo.EllipticCurve == entities.Secp256k1:
		return s.signECDSA(privkey, data)
	default:
		errMessage := "signing algorithm and curve combination not supported for signing"
		logger.With("algorithm", key.Algo.Type, "curve", key.Algo.EllipticCurve).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
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

func (s *Store) createKey(ctx context.Context, id string, importedPrivKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("creating key")

	var privKey []byte
	var pubKey []byte
	switch {
	case alg.Type == entities.Eddsa && alg.EllipticCurve == entities.Bn254:
		eddsaKey, err := eddsaBN254(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate EDDSA/BN254 key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		privKey = eddsaKey.Bytes()
		pubKey = eddsaKey.Public().Bytes()
	case alg.Type == entities.Ecdsa && alg.EllipticCurve == entities.Secp256k1:
		ecdsaKey, err := ecdsaSecp256k1(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate Secp256k1/ECDSA key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		privKey = crypto.FromECDSA(ecdsaKey)
		pubKey = crypto.FromECDSAPub(&ecdsaKey.PublicKey)
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	key := &entities.Key{
		ID:        id,
		PublicKey: pubKey,
		Algo:      alg,
		Tags:      attr.Tags,
		Metadata:  &entities.Metadata{},
	}
	err := s.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		err := dbtx.Keys().Add(ctx, models.NewKey(key))
		if err != nil {
			return err
		}

		_, err = s.secretStore.Set(ctx, id, base64.StdEncoding.EncodeToString(privKey), attr)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("key was created successfully")
	return key, nil
}

func (s *Store) signECDSA(privKey, data []byte) ([]byte, error) {
	if len(data) != crypto.DigestLength {
		return nil, errors.InvalidParameterError("data is required to be exactly %d bytes (%d)", crypto.DigestLength, len(data))
	}

	ecdsaPrivKey, err := crypto.ToECDSA(privKey)
	if err != nil {
		errMessage := "failed to parse ECDSA private key"
		s.logger.With("error", err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signature, err := crypto.Sign(data, ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign payload with ECDSA"
		s.logger.With("error", err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	s.logger.Debug("payload signed successfully")
	// We remove the recID from the signature (last byte).
	return signature[:len(signature)-1], nil
}

func (s *Store) signEDDSA(privKeyB, data []byte) ([]byte, error) {
	privKey := eddsabn254.PrivateKey{}
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		errMessage := "failed to parse EDDSA private key"
		s.logger.With("error", err).Error(errMessage)
		return nil, errors.DependencyFailureError(errMessage)
	}

	signature, err := privKey.Sign(data, hash.MIMC_BN254.New("seed"))
	if err != nil {
		errMessage := "failed to sign payload with EDDSA"
		s.logger.With("error", err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	s.logger.Debug("payload signed successfully")
	return signature, nil
}

func eddsaBN254(importedPrivKey []byte) (eddsabn254.PrivateKey, error) {
	if importedPrivKey == nil {
		seed := make([]byte, 32)
		rand.New(rand.NewSource(time.Now().UnixNano())).Read(seed)

		// Usually standards implementations of eddsa do not require the choice of a specific hash function (usually it's SHA256).
		// Here we needed to allow the choice of the hash so we can chose a hash function that is easily programmable in a snark circuit.
		// Same hFunc should be used for sign and verify
		return eddsabn254.GenerateKey(bytes.NewReader(seed))
	}

	key := eddsabn254.PrivateKey{}
	_, err := key.SetBytes(importedPrivKey)
	if err != nil {
		return key, err
	}

	return key, nil
}

func ecdsaSecp256k1(importedPrivKey []byte) (*ecdsa.PrivateKey, error) {
	if importedPrivKey == nil {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	return crypto.ToECDSA(importedPrivKey)
}
