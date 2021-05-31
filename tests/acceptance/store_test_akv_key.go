package acceptancetests

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/keys/akv"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type akvKeyTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *akv.Store
}

func (s *akvKeyTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.Run("should create a new key pair successfully", func() {
		id := fmt.Sprintf("my-key-create-%d", common.RandInt(1000))
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.NotNil(s.T(), key.PublicKey)
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.Version)
		assert.NotNil(s.T(), key.Metadata.CreatedAt)
		assert.NotNil(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)

		err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), key)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})
}

func (s *akvKeyTestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.Run("should import a new key pair successfully: ECDSA/Secp256k1", func() {
		id := fmt.Sprintf("my-key-ecdsa-import-%d", common.RandInt(1000))
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.Version)
		assert.NotNil(s.T(), key.Metadata.CreatedAt)
		assert.NotNil(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)

		err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

	s.Run("should fail to import a new key pair: EDDSA/BN254 (not implemented yet)", func() {
		id := "my-key-eddsa-import"
		privKey, _ := hex.DecodeString(privKeyEDDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), key)
		assert.Equal(s.T(), err, errors.ErrNotSupported)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), key)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})
}

func (s *akvKeyTestSuite) TestGet() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-key-get-%d", common.RandInt(1000))
	tags := testutils.FakeTags()
	privKey, _ := hex.DecodeString(privKeyECDSA)

	key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	defer func() {
		err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.Run("should get a key pair successfully", func() {
		keyRetrieved, err := s.store.Get(ctx, id)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, keyRetrieved.ID)
		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), tags, keyRetrieved.Tags)
		assert.Equal(s.T(), entities.Secp256k1, keyRetrieved.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, keyRetrieved.Algo.Type)
		assert.NotEmpty(s.T(), keyRetrieved.Metadata.Version)
		assert.NotNil(s.T(), keyRetrieved.Metadata.CreatedAt)
		assert.NotNil(s.T(), keyRetrieved.Metadata.UpdatedAt)
		assert.True(s.T(), keyRetrieved.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), keyRetrieved.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), keyRetrieved.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), keyRetrieved.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		keyRetrieved, getErr := s.store.Get(ctx, "invalidID")

		require.Nil(s.T(), keyRetrieved)
		assert.True(s.T(), errors.IsNotFoundError(getErr))
	})
}

func (s *akvKeyTestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := fmt.Sprintf("my-key-list-%s", common.RandString(5))

	_, err := s.store.Create(ctx, id, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	defer func() {
		err := s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.Run("should list all key pairs", func() {
		ids, err := s.store.List(ctx)
		require.NoError(s.T(), err)
		assert.Contains(s.T(), ids, id)
	})
}

func (s *akvKeyTestSuite) TestSign() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	payload := crypto.Keccak256([]byte("my data to sign"))
	privKey, _ := hex.DecodeString(privKeyECDSA)

	id := fmt.Sprintf("mykey-sign-ecdsa-%d", common.RandInt(1000))
	_, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	defer func() {
		err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.Run("should sign a message successfully: ECDSA/Secp256k1", func() {
		signature, err := s.store.Sign(ctx, id, payload)
		require.NoError(s.T(), err)

		verified, err := verifySignature(signature, payload, privKey)
		require.NoError(s.T(), err)
		require.True(s.T(), verified)
	})

	s.Run("should fail and parse the error code correctly", func() {
		signature, signErr := s.store.Sign(ctx, "invalidID", payload)

		require.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsNotFoundError(signErr))
	})
}

func verifySignature(signature, msg, privKeyB []byte) (bool, error) {
	privKey, err := crypto.ToECDSA(privKeyB)
	if err != nil {
		return false, err
	}
	fmt.Println(privKey.PublicKey.X.String(), privKey.PublicKey.Y.String())

	if len(signature) == EthSignatureLength {
		retrievedPubkey, err := crypto.SigToPub(msg, signature)
		if err != nil {
			return false, err
		}

		fmt.Println(retrievedPubkey.X.String(), retrievedPubkey.Y.String())

		return privKey.PublicKey.Equal(retrievedPubkey), nil
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])
	return ecdsa.Verify(&privKey.PublicKey, msg, r, s), nil
}
