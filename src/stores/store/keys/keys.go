package keys

import (
	"context"
	"encoding/base64"

	"github.com/consensysquorum/quorum-key-manager/pkg/crypto"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=keys.go -destination=mock/keys.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Create a new key and stores it
	Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error)

	// Import an externally created key and stores it
	Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error)

	// Get the public part of a stored key.
	Get(ctx context.Context, id string) (*entities.Key, error)

	// List keys
	List(ctx context.Context) ([]string, error)

	// Update key tags
	Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error)

	// Delete secret not permanently, by using Undelete() the secret can be retrieve
	Delete(ctx context.Context, id string) error

	// GetDeleted keys
	GetDeleted(ctx context.Context, id string) (*entities.Key, error)

	// ListDeleted keys
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted secret
	Undelete(ctx context.Context, id string) error

	// Destroy secret permanently
	Destroy(ctx context.Context, id string) error

	// Sign from any arbitrary data using the specified key
	Sign(ctx context.Context, id string, data []byte) ([]byte, error)

	// Verify verifies the signature belongs to the corresponding key
	Verify(ctx context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error

	// Encrypt any arbitrary data using a specified key
	Encrypt(ctx context.Context, id string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, id string, data []byte) ([]byte, error)
}

func VerifySignature(logger log.Logger, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	logger = logger.With(
		"pub_key", base64.URLEncoding.EncodeToString(pubKey),
		"curve", algo.EllipticCurve,
		"signing_algorithm", algo.Type,
	)

	var err error
	var verified bool
	switch {
	case algo.EllipticCurve == entities.Secp256k1 && algo.Type == entities.Ecdsa:
		verified, err = crypto.VerifyECDSASignature(pubKey, data, sig)
	case algo.EllipticCurve == entities.Bn254 && algo.Type == entities.Eddsa:
		verified, err = crypto.VerifyEDDSASignature(pubKey, data, sig)
	default:
		errMessage := "unsupported algorithm"
		logger.Error(errMessage)
		return errors.NotSupportedError(errMessage)
	}
	if err != nil {
		errMessage := "failed to verify signature"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	if !verified {
		errMessage := "signature does not belong to the specified public key"
		logger.Error(errMessage, "pub_key", pubKey)
		return errors.InvalidParameterError(errMessage)
	}

	return nil
}
