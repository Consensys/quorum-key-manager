package utils

import (
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/crypto/ecdsa"
	"github.com/consensys/quorum-key-manager/pkg/crypto/eddsa"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	invalidPublicKey = []byte("invalid pub key")
	invalidSignature = []byte("invalid signature")
)

func TestKeysVerifyMessage_ecdsa256k1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	privKey, pubKey, _ := ecdsa.CreateSecp256k1(nil)
	_, pubKey2, _ := ecdsa.CreateSecp256k1(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := ecdsa.SignSecp256k1(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := ecdsa.SignSecp256k1(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding public key", func(t *testing.T) {
		err := connector.Verify(pubKey2, data, signature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding key type", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail verify invalid public key size", func(t *testing.T) {
		err := connector.Verify(invalidPublicKey, data, signature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail verify invalid signature format", func(t *testing.T) {
		err := connector.Verify(pubKey, data, invalidSignature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func TestKeysVerifyMessage_eddsaBabyJubJub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	privKey, pubKey, _ := eddsa.CreateBabyjubjub(nil)
	_, pubKey2, _ := eddsa.CreateBabyjubjub(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := eddsa.SignBabyjubjub(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := eddsa.SignBabyjubjub(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding public key", func(t *testing.T) {
		err := connector.Verify(pubKey2, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding key type", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail verify invalid public key size", func(t *testing.T) {
		err := connector.Verify(invalidPublicKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail verify invalid signature format", func(t *testing.T) {
		err := connector.Verify(pubKey, data, invalidSignature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func TestKeysVerifyMessage_ed25519(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	privKey, pubKey, _ := eddsa.CreateX25519(nil)
	_, pubKey2, _ := eddsa.CreateX25519(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := eddsa.SignX25519(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := eddsa.SignX25519(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding public key", func(t *testing.T) {
		err := connector.Verify(pubKey2, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding key type", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify invalid signature size", func(t *testing.T) {
		err := connector.Verify(pubKey, data, invalidSignature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify invalid public key size", func(t *testing.T) {
		err := connector.Verify(invalidPublicKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func TestKeysVerifyMessage_errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)
	connector := New(logger)
	t.Run("should fail to not support types", func(t *testing.T) {
		err := connector.Verify(nil, nil, nil, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsNotSupportedError(err))
	})
}
