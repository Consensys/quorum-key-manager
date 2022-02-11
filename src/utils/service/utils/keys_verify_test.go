package utils

import (
	"testing"

	pkgcrypto "github.com/consensys/quorum-key-manager/pkg/crypto"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeysVerifyMessage_ecdsa256k1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	privKey, pubKey, _ := pkgcrypto.ECDSASecp256k1(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := pkgcrypto.SignECDSA256k1(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := pkgcrypto.SignECDSA256k1(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail to verify no corresponding signing algo", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
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
	privKey, pubKey, _ := pkgcrypto.EdDSABabyjubjub(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := pkgcrypto.SignEDDSABabyjubjub(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := pkgcrypto.SignEDDSABabyjubjub(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
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
	privKey, pubKey, _ := pkgcrypto.EdDSA25519(nil)
	data := crypto.Keccak256([]byte("my data to sign"))
	signature, err := pkgcrypto.SignEDDSA25519(privKey, data)
	require.NoError(t, err)

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		assert.NoError(t, err)
	})

	t.Run("should fail to verify no corresponding signature", func(t *testing.T) {
		invalidSig, _ := pkgcrypto.SignEDDSA25519(privKey, crypto.Keccak256([]byte("invalid data")))
		err := connector.Verify(pubKey, data, invalidSig, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(pubKey, data, signature, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Babyjubjub,
		})

		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func TestKeysVerifyMessage_notSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := testutils.NewMockLogger(ctrl)

	connector := New(logger)
	t.Run("should verify message successfully", func(t *testing.T) {
		err := connector.Verify(nil, nil, nil, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.X25519,
		})

		require.Error(t, err)
		assert.True(t, errors.IsNotSupportedError(err))
	})
}
