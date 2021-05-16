package acceptancetests

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/akv"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	s.T().Run("should create a new key pair successfully", func(t *testing.T) {
		id := fmt.Sprintf("my-key-create-%d", common.RandInt(1000))
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, key.ID)
		assert.NotNil(t, key.PublicKey)
		assert.Equal(t, tags, key.Tags)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.NotEmpty(t, key.Metadata.Version)
		assert.NotNil(t, key.Metadata.CreatedAt)
		assert.NotNil(t, key.Metadata.UpdatedAt)
		assert.True(t, key.Metadata.DeletedAt.IsZero())
		assert.True(t, key.Metadata.DestroyedAt.IsZero())
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.False(t, key.Metadata.Disabled)

		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		id := "my-key"
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *akvKeyTestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.T().Run("should import a new key pair successfully: ECDSA/Secp256k1", func(t *testing.T) {
		id := fmt.Sprintf("my-key-ecdsa-import-%d", common.RandInt(1000))

		key, err := s.store.Import(ctx, id, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, key.ID)
		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", key.PublicKey)
		assert.Equal(t, tags, key.Tags)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.NotEmpty(t, key.Metadata.Version)
		assert.NotNil(t, key.Metadata.CreatedAt)
		assert.NotNil(t, key.Metadata.UpdatedAt)
		assert.True(t, key.Metadata.DeletedAt.IsZero())
		assert.True(t, key.Metadata.DestroyedAt.IsZero())
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.False(t, key.Metadata.Disabled)

		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	})

	s.T().Run("should fail to import a new key pair: EDDSA/BN254 (not implemented yet)", func(t *testing.T) {
		id := "my-key-eddsa-import"
		tags := testutils.FakeTags()

		key, err := s.store.Import(ctx, id, "5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7", &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, key)
		assert.Equal(t, err, errors.ErrNotSupported)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		id := "my-key"
		tags := testutils.FakeTags()

		key, err := s.store.Import(ctx, id, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *akvKeyTestSuite) TestGet() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-key-get-%d", common.RandInt(1000))
	tags := testutils.FakeTags()

	key, err := s.store.Import(ctx, id, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	defer func() {
		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.T().Run("should get a key pair successfully", func(t *testing.T) {
		keyRetrieved, err := s.store.Get(ctx, id, "")
		require.NoError(t, err)

		assert.Equal(t, id, keyRetrieved.ID)
		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", key.PublicKey)
		assert.Equal(t, tags, keyRetrieved.Tags)
		assert.Equal(t, entities.Secp256k1, keyRetrieved.Algo.EllipticCurve)
		assert.Equal(t, entities.Ecdsa, keyRetrieved.Algo.Type)
		assert.NotEmpty(t, keyRetrieved.Metadata.Version)
		assert.NotNil(t, keyRetrieved.Metadata.CreatedAt)
		assert.NotNil(t, keyRetrieved.Metadata.UpdatedAt)
		assert.True(t, keyRetrieved.Metadata.DeletedAt.IsZero())
		assert.True(t, keyRetrieved.Metadata.DestroyedAt.IsZero())
		assert.True(t, keyRetrieved.Metadata.ExpireAt.IsZero())
		assert.False(t, keyRetrieved.Metadata.Disabled)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		keyRetrieved, err := s.store.Get(ctx, "invalidID", "")

		require.Nil(t, keyRetrieved)
		assert.True(t, errors.IsNotFoundError(err))
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
		_, err := s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.T().Run("should list all key pairs", func(t *testing.T) {
		keys, err := s.store.List(ctx)
		require.NoError(t, err)

		assert.NotEmpty(t, keys)
	})
}

func (s *akvKeyTestSuite) TestSign() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	payload := hexutil.Encode([]byte("my data to sign"))
	privKey := "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"

	id := fmt.Sprintf("mykey-sign-ecdsa-%d", common.RandInt(1000))
	_, err := s.store.Import(ctx, id, privKey, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	defer func() {
		_, err = s.store.Delete(ctx, id)
		require.NoError(s.T(), err)
		_ = s.store.Destroy(ctx, id)
	}()

	s.T().Run("should sign a message successfully: ECDSA/Secp256k1", func(t *testing.T) {
		signature, err := s.store.Sign(ctx, id, payload, "")
		require.NoError(t, err)

		verified, err := verifySignature(signature, payload, privKey)
		require.NoError(t, err)
		require.True(t, verified)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		id := "my-key"

		key, err := s.store.Sign(ctx, id, "", "")

		require.Empty(t, key)
		assert.True(t, errors.IsInvalidFormatError(err))
	})
}

func verifySignature(signature, msg, privKey string) (bool, error) {
	bSig, _ := hexutil.Decode(signature)
	bMsg, _ := hexutil.Decode(msg)
	privKeyS, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return false, err
	}

	r := new(big.Int).SetBytes(bSig[0:32])
	s := new(big.Int).SetBytes(bSig[32:64])

	return ecdsa.Verify(&privKeyS.PublicKey, crypto.Keccak256(bMsg), r, s), nil
}