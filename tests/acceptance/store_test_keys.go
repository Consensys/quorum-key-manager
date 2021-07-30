package acceptancetests

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/connectors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
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
	env       *IntegrationEnvironment
	connector connectors.KeysConnector
}

func (s *keysTestSuite) TearDownSuite() {
	ctx := s.env.ctx

	keyIds, err := s.connector.List(ctx)
	require.NoError(s.T(), err)
	s.env.logger.Info("deleting the following keys", "keys", keyIds)

	for _, id := range keyIds {
		err = s.connector.Delete(ctx, id)
		require.NoError(s.T(), err)
	}

	for _, id := range keyIds {
		err = s.connector.Destroy(ctx, id)
		if err == nil {
			break
		}
	}
}

func (s *keysTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.Run("should create a new key pair successfully", func() {
		id := s.newID("my-key-create")
		tags := testutils.FakeTags()

		key, err := s.connector.Create(ctx, id, &entities.Algorithm{
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
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		tags := testutils.FakeTags()

		key, err := s.connector.Create(ctx, id, &entities.Algorithm{
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
		id := fmt.Sprintf("%s-%d", "my-key-ecdsa-import", common.RandInt(10000))
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.connector.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
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
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should import a new key pair successfully: EDDSA/BN254", func() {
		id := fmt.Sprintf("%s-%d", "my-key-eddsa-import", common.RandInt(10000))
		privKey, _ := hex.DecodeString(privKeyEDDSA)

		key, err := s.connector.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
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
		assert.Equal(s.T(), entities.Bn254, key.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Eddsa, key.Algo.Type)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), key.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), key.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		id := "my-key"
		privKey, _ := hex.DecodeString(privKeyECDSA)

		key, err := s.connector.Import(ctx, id, privKey, &entities.Algorithm{
			Type:          entities.Ecdsa,
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

	_, err := s.connector.Create(ctx, id, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get a key pair successfully", func() {
		keyRetrieved, err := s.connector.Get(ctx, id)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, keyRetrieved.ID)
		assert.Equal(s.T(), tags, keyRetrieved.Tags)
		assert.Equal(s.T(), entities.Secp256k1, keyRetrieved.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, keyRetrieved.Algo.Type)
		assert.NotEmpty(s.T(), keyRetrieved.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), keyRetrieved.Metadata.UpdatedAt)
		assert.True(s.T(), keyRetrieved.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), keyRetrieved.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), keyRetrieved.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), keyRetrieved.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		keyRetrieved, getErr := s.connector.Get(ctx, "invalidID")

		require.Nil(s.T(), keyRetrieved)
		assert.True(s.T(), errors.IsNotFoundError(getErr))
	})
}

func (s *keysTestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := s.newID("my-key-list")

	_, err := s.connector.Create(ctx, id, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should list all key pairs", func() {
		ids, err := s.connector.List(ctx)
		require.NoError(s.T(), err)
		assert.Contains(s.T(), ids, id)
	})
}

func (s *keysTestSuite) TestUpdate() {
	ctx := s.env.ctx
	id := s.newID("my-key-update")
	tags := testutils.FakeTags()
	_, err := s.connector.Create(ctx, id, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
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

		updatedKey, err := s.connector.Update(ctx, id, &entities.Attributes{
			Tags: newTags,
		})

		require.NoError(s.T(), err)

		assert.Equal(s.T(), id, updatedKey.ID)
		assert.Equal(s.T(), newTags, updatedKey.Tags)
		assert.Equal(s.T(), entities.Secp256k1, updatedKey.Algo.EllipticCurve)
		assert.Equal(s.T(), entities.Ecdsa, updatedKey.Algo.Type)
		assert.NotEmpty(s.T(), updatedKey.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), updatedKey.Metadata.UpdatedAt)
		assert.True(s.T(), updatedKey.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), updatedKey.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), updatedKey.Metadata.ExpireAt.IsZero())
		assert.False(s.T(), updatedKey.Metadata.Disabled)
	})

	s.Run("should fail and parse the error code correctly", func() {
		updatedKey, err := s.connector.Update(ctx, "invalidID", &entities.Attributes{
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
		key, err := s.connector.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		signature, err := s.connector.Sign(ctx, id, payload)
		require.NoError(s.T(), err)

		err = s.connector.Verify(ctx, key.PublicKey, payload, signature, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should sign and verify a message successfully: EDDSA/BN254", func() {
		id := fmt.Sprintf("mykey-sign-eddsa-%d", common.RandInt(1000))
		payload := []byte("my data to sign")
		key, err := s.connector.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})
		if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsInvalidParameterError(err) {
			return
		}
		require.NoError(s.T(), err)

		signature, err := s.connector.Sign(ctx, id, payload)
		require.NoError(s.T(), err)

		err = s.connector.Verify(ctx, key.PublicKey, payload, signature, &entities.Algorithm{
			Type:          key.Algo.Type,
			EllipticCurve: key.Algo.EllipticCurve,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should fail and parse the error code correctly", func() {
		signature, signErr := s.connector.Sign(ctx, "invalidID", crypto.Keccak256([]byte("my data to sign")))

		require.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsNotFoundError(signErr))
	})
}

func (s *keysTestSuite) newID(name string) string {
	return fmt.Sprintf("%s-%s", name, common.RandHexString(16))
}
