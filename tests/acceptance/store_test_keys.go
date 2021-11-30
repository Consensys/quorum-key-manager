package acceptancetests

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privKeyECDSA = "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
	privKeyEDDSA = "5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7"
)

type keysTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store stores.KeyStore
	utils stores.Utilities
	db    database.Keys
}

func (s *keysTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.Run("should create a new key pair successfully", func() {
		id := s.newID("my-key-create")
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.NotNil(s.T(), key.PublicKey)
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities2.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should create a new key pair successfully if it already exists in the Vault", func() {
		id := s.newID("my-key-create")
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		err = s.db.Delete(ctx, id)
		require.NoError(s.T(), err)
		err = s.db.Purge(ctx, id)
		require.NoError(s.T(), err)

		key, err = s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.NotNil(s.T(), key.PublicKey)
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities2.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		tags := testutils.FakeTags()

		key, err := s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
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
		id := fmt.Sprintf("%s-%d", "my-key-ecdsa-import", common.RandInt(10000))
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		if err != nil && errors.IsNotSupportedError(err) {
			return
		}
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities2.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should import a new key pair successfully: EDDSA/Babyjubjub", func() {
		id := fmt.Sprintf("%s-%d", "my-key-eddsa-import", common.RandInt(10000))
		privKey, _ := hex.DecodeString(privKeyEDDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities2.Algorithm{
			Type:          entities2.Eddsa,
			EllipticCurve: entities2.Babyjubjub,
		}, &entities.Attributes{
			Tags: tags,
		})
		// AKV and AWS does not support EDDSA
		if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsInvalidParameterError(err) {
			assert.Nil(s.T(), key)
			return
		}
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), tags, key.Tags)
		assert.Equal(s.T(), entities2.Babyjubjub, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Eddsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.store.Import(ctx, id, privKey, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})
		if err != nil && errors.IsNotSupportedError(err) {
			return
		}

		require.Nil(s.T(), key)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})
}

func (s *keysTestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-key-get")
	tags := testutils.FakeTags()

	_, err := s.store.Create(ctx, id, &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get a key pair successfully", func() {
		keyRetrieved, err := s.store.Get(ctx, id)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, keyRetrieved.ID)
		assert.Equal(s.T(), tags, keyRetrieved.Tags)
		assert.Equal(s.T(), entities2.Secp256k1, keyRetrieved.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Ecdsa, keyRetrieved.Algo.Type)
		assert.NotEmpty(s.T(), keyRetrieved.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), keyRetrieved.Metadata.UpdatedAt)
		assert.True(s.T(), keyRetrieved.Metadata.DeletedAt.IsZero())
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
	id2 := s.newID("my-key-list-2")
	id3 := s.newID("my-key-list-3")

	_, err := s.store.Create(ctx, id, &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	_, err = s.store.Create(ctx, id2, &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	_, err = s.store.Create(ctx, id3, &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	listLen := 0
	s.Run("should list all key pairs", func() {
		ids, err := s.store.List(ctx, 0, 0)
		require.NoError(s.T(), err)

		listLen = len(ids)
		assert.Contains(s.T(), ids, id)
		assert.Contains(s.T(), ids, id2)
		assert.Contains(s.T(), ids, id3)
	})

	s.Run("should list first key pair successfully", func() {
		ids, err := s.store.List(ctx, 1, uint64(listLen-3))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id})
	})

	s.Run("should list last two key pair successfully", func() {
		ids, err := s.store.List(ctx, 2, uint64(listLen-2))

		require.NoError(s.T(), err)
		assert.Equal(s.T(), ids, []string{id2, id3})
	})
}

func (s *keysTestSuite) TestUpdate() {
	ctx := s.env.ctx
	id := s.newID("my-key-update")
	tags := testutils.FakeTags()
	_, err := s.store.Create(ctx, id, &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)
	require.NoError(s.T(), err)

	s.Run("should update a key pair successfully", func() {
		newTags := map[string]string{
			"newTag1": "tagValue1",
			"newTag2": "tagValue2",
		}

		updatedKey, err := s.store.Update(ctx, id, &entities.Attributes{
			Tags: newTags,
		})

		require.NoError(s.T(), err)
		require.NotNil(s.T(), updatedKey)

		assert.Equal(s.T(), id, updatedKey.ID)
		assert.Equal(s.T(), newTags, updatedKey.Tags)
		assert.Equal(s.T(), entities2.Secp256k1, updatedKey.Algo.EllipticCurve)
		assert.Equal(s.T(), entities2.Ecdsa, updatedKey.Algo.Type)
		assert.NotEmpty(s.T(), updatedKey.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), updatedKey.Metadata.UpdatedAt)
		assert.True(s.T(), updatedKey.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), updatedKey.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), updatedKey.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		updatedKey, err := s.store.Update(ctx, "invalidID", &entities.Attributes{
			Tags: testutils.FakeTags(),
		})

		require.Nil(s.T(), updatedKey)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *keysTestSuite) TestSignVerify() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.Run("should sign and verify a message successfully: ECDSA/Secp256k1", func() {
		id := s.newID("mykey-sign-ecdsa")
		payload := crypto.Keccak256([]byte("my data to sign"))
		key, err := s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		signature, err := s.store.Sign(ctx, id, payload, nil)
		require.NoError(s.T(), err)

		err = s.utils.Verify(key.PublicKey, payload, signature, &entities2.Algorithm{
			Type:          entities2.Ecdsa,
			EllipticCurve: entities2.Secp256k1,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should sign and verify a message successfully: EDDSA/Babyjubjub", func() {
		id := fmt.Sprintf("mykey-sign-eddsa-%d", common.RandInt(1000))
		payload := []byte("my data to sign")
		key, err := s.store.Create(ctx, id, &entities2.Algorithm{
			Type:          entities2.Eddsa,
			EllipticCurve: entities2.Babyjubjub,
		}, &entities.Attributes{
			Tags: tags,
		})
		if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsInvalidParameterError(err) {
			return
		}
		require.NoError(s.T(), err)

		signature, err := s.store.Sign(ctx, id, payload, nil)
		require.NoError(s.T(), err)

		err = s.utils.Verify(key.PublicKey, payload, signature, &entities2.Algorithm{
			Type:          key.Algo.Type,
			EllipticCurve: key.Algo.EllipticCurve,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should fail and parse the error code correctly", func() {
		signature, signErr := s.store.Sign(ctx, "invalidID", crypto.Keccak256([]byte("my data to sign")), nil)

		require.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsNotFoundError(signErr))
	})
}

func (s *keysTestSuite) newID(name string) string {
	return fmt.Sprintf("%s-%s", name, common.RandHexString(16))
}
