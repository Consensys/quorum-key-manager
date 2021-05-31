package acceptancetests

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privKeyECDSA       = "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
	privKeyECDSA2      = "5a1e076fd6b1a0daf31fd1cc0b525ea230f9e50d06f002daff271315262f06fa"
	privKeyEDDSA       = "5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7"
	EthSignatureLength = 65
)

type keysTestSuite struct {
	suite.Suite
	env    *IntegrationEnvironment
	store  keys.Store
	keyIds []string
}

func (s *keysTestSuite) TearDownSuite() {
	ctx := s.env.ctx

	s.env.logger.WithField("keys", s.keyIds).Info("Deleting the following keys")
	for _, id := range s.keyIds {
		err := s.store.Delete(ctx, id)
		if err != nil && errors.IsNotSupportedError(err) {
			return
		}
	}

	for _, id := range s.keyIds {
		_ = s.store.Destroy(ctx, id)
	}
}

func (s *keysTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.Run("should create a new key pair successfully", func() {
		id := s.newID("my-key-create")
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

func (s *keysTestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.Run("should import a new key pair successfully: ECDSA/Secp256k1", func() {
		id := s.newID("my-key-ecdsa-import")
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
	})

	s.Run("should import a new key pair successfully: EDDSA/BN254", func() {
		id := "my-key-eddsa-import"
		privKey, _ := hex.DecodeString(privKeyEDDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})
		// AKV and AWS does not support EDDSA
		if err != nil && errors.IsNotSupportedError(err) {
			assert.Nil(s.T(), key)
			return
		}

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities.Bn254, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Eddsa, key.Algo.Type)
		assert.Equal(s.T(), "1", key.Metadata.Version)
		assert.NotNil(s.T(), key.Metadata.CreatedAt)
		assert.NotNil(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
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

func (s *keysTestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-key-get")
	tags := testutils.FakeTags()
	privKey, _ := hex.DecodeString(privKeyECDSA)

	key, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

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

func (s *keysTestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := s.newID("my-key-list")

	_, err := s.store.Create(ctx, id, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should list all key pairs", func() {
		ids, err := s.store.List(ctx)
		require.NoError(s.T(), err)
		assert.Contains(s.T(), ids, id)
	})
}

func (s *keysTestSuite) TestSign() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	payload := crypto.Keccak256([]byte("my data to sign"))
	privKey, _ := hex.DecodeString(privKeyECDSA)
	id := s.newID("mykey-sign-ecdsa")

	_, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

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

	if len(signature) == EthSignatureLength {
		retrievedPubkey, err := crypto.SigToPub(crypto.Keccak256(msg), signature)
		if err != nil {
			return false, err
		}

		return privKey.PublicKey.Equal(retrievedPubkey), nil
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])
	return ecdsa.Verify(&privKey.PublicKey, msg, r, s), nil
}

func (s *keysTestSuite) newID(name string) string {
	id := fmt.Sprintf("%s-%d", name, common.RandInt(1000))
	s.keyIds = append(s.keyIds, id)

	return id
}
