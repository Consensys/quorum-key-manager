// +build acceptance

package integrationtests

import "C"
import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/akv"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type akvKeyTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *akv.KeyStore
}

func (s *akvKeyTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.T().Run("should create a new key pair successfully", func(t *testing.T) {
		id := "my-key-create"
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
		id := "my-key-ecdsa-import"

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
	})

	s.T().Run("should import a new key pair successfully: EDDSA/BN254", func(t *testing.T) {
		id := "my-key-eddsa-import"
		tags := testutils.FakeTags()

		key, err := s.store.Import(ctx, id, "5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7", &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, key)
		assert.Equal(t, err, errors.NotImplementedError)
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
	id := "my-key-get"
	tags := testutils.FakeTags()

	key, err := s.store.Import(ctx, id, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

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
}

func (s *akvKeyTestSuite) TestList() {
	ctx := s.env.ctx
	id1 := "my-key-list1"
	tags := testutils.FakeTags()

	_, err := s.store.Import(ctx, id1, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.T().Run("should list all key pairs", func(t *testing.T) {
		keys, err := s.store.List(ctx)
		require.NoError(t, err)

		assert.Contains(t, keys, id1)
	})
}

func (s *akvKeyTestSuite) TestSign() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	payload := hexutil.Encode([]byte("my data to sign"))

	s.T().Run("should sign a message successfully: ECDSA/Secp256k1", func(t *testing.T) {
		id := fmt.Sprintf("mykeysignecdsa-%d", rand.Intn(1000))
		iKey, err := s.store.Import(ctx, id, "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c", &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(t, err)
		

		signature, err := s.store.Sign(ctx, id, payload, iKey.Metadata.Version)
		require.NoError(t, err)
		assert.True(t, verifySignature(signature, payload, iKey.PublicKey))
		assert.Equal(t, "0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed01", signature)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		id := "my-key"
	
		key, err := s.store.Sign(ctx, id, "", "")
	
		require.Empty(t, key)
		assert.True(t, errors.IsInvalidFormatError(err))
	})
}

func verifySignature(signature, msg, pubkey string) bool {
	bSig, _ := hexutil.Decode(signature)
	bMsg, _ := hexutil.Decode(msg)
	bPubKey, _ := hexutil.Decode(pubkey)
	verified := secp256k1.VerifySignature(bSig, bMsg, bPubKey)
	return verified
}

