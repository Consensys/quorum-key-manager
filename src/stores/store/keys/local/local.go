package local

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/gnark-crypto/crypto/hash"
	eddsabn254 "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

type Store struct {
	secretStore secrets.Store
	logger      log.Logger
}

var _ keys.Store = &Store{}

func New(secretStore secrets.Store, logger log.Logger) *Store {
	return &Store{
		secretStore: secretStore,
		logger:      logger,
	}
}

func (s *Store) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.create(ctx, id, nil, alg, attr)
}

func (s *Store) Import(ctx context.Context, id string, importedPrivKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	return s.create(ctx, id, importedPrivKey, alg, attr)
}

func (s *Store) create(ctx context.Context, id string, importedPrivKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.With("id", id).With("signing_algorithm", alg.Type).With("curve", alg.EllipticCurve)

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

	secret, err := s.secretStore.Set(ctx, id, base64.StdEncoding.EncodeToString(privKey), &entities.Attributes{})
	if err != nil {
		return nil, err
	}

	return &entities.Key{
		ID:        id,
		PublicKey: pubKey,
		Algo: &entities.Algorithm{
			Type:          alg.Type,
			EllipticCurve: alg.EllipticCurve,
		},
		Metadata: &entities.Metadata{
			Disabled:  false,
			CreatedAt: secret.Metadata.CreatedAt,
			UpdatedAt: secret.Metadata.UpdatedAt,
		},
		Tags: attr.Tags,
	}, nil
}

func (s *Store) Update(_ context.Context, _ string, _ *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotSupported
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.secretStore.Delete(ctx, id)
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	return s.secretStore.Undelete(ctx, id)
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	return s.secretStore.Destroy(ctx, id)
}

func (s *Store) Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error) {
	logger := s.logger.With("id", id)

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
	case algo.Type == entities.Eddsa && algo.EllipticCurve == entities.Bn254:
		return s.signEDDSA(privkey, data)
	case algo.Type == entities.Ecdsa && algo.EllipticCurve == entities.Secp256k1:
		return s.signECDSA(privkey, data)
	default:
		errMessage := "signing algorithm and curve combination not supported for signing"
		logger.With("algorithm", algo.Type, "curve", algo.EllipticCurve).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}

func (s *Store) Encrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Decrypt(_ context.Context, id string, data []byte) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) signECDSA(privKey, data []byte) ([]byte, error) {
	if len(data) != crypto.DigestLength {
		errMessage := fmt.Sprintf("data is required to be exactly %d bytes (%d)", crypto.DigestLength, len(data))
		s.logger.With("data_length", len(data), "expected_data_length", crypto.DigestLength).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
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

	return signature, nil
}

func eddsaBN254(importedPrivKey []byte) (eddsabn254.PrivateKey, error) {
	if importedPrivKey == nil {
		seed := make([]byte, 32)
		rand.New(rand.NewSource(time.Now().UnixNano())).Read(seed)

		// Usually standards implementations of eddsa do not require the choice of a specific hash function (usually it's SHA256).
		// Here we needed to allow the choice of the hash, so we can choose a hash function that is easily programmable in a snark circuit.
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
