package acceptancetests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/hashicorp"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TODO: Destroy key pairs when done with the tests to avoid conflicts between tests

const (
	ecdsaPrivKeyb64 = "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="
	eddsaPrivKeyb64 = "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZGGbioek5qYuzJzTNZpTHrVjjFk7iFe3FYwfpxZyNPxtIaFB5gb9VP9IcHZewwNZly821re7RkmB8pGdjywygPH"
)

type hashicorpKeyTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store *hashicorp.Store
}

func (s *hashicorpKeyTestSuite) TestCreate() {
	ctx := s.env.ctx

	s.T().Run("should create a new key pair successfully", func(t *testing.T) {
		id := fmt.Sprintf("my-key-create-%d", rand.Intn(1000))
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
		assert.Equal(t, "1", key.Metadata.Version)
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

func (s *hashicorpKeyTestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.T().Run("should import a new key pair successfully: ECDSA/Secp256k1", func(t *testing.T) {
		id := fmt.Sprintf("my-key-ecdsa-import-%d", rand.Intn(1000))

		key, err := s.store.Import(ctx, id, ecdsaPrivKeyb64, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, key.ID)
		assert.Equal(t, "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", key.PublicKey)
		assert.Equal(t, tags, key.Tags)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, "1", key.Metadata.Version)
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

		key, err := s.store.Import(ctx, id, eddsaPrivKeyb64, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)

		assert.Equal(t, id, key.ID)
		assert.Equal(t, "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", key.PublicKey)
		assert.Equal(t, tags, key.Tags)
		assert.Equal(t, entities.Bn254, key.Algo.EllipticCurve)
		assert.Equal(t, entities.Eddsa, key.Algo.Type)
		assert.Equal(t, "1", key.Metadata.Version)
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

		key, err := s.store.Import(ctx, id, ecdsaPrivKeyb64, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: "invalidCurve",
		}, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *hashicorpKeyTestSuite) TestGet() {
	ctx := s.env.ctx
	id := "my-key-get"
	tags := testutils.FakeTags()

	key, err := s.store.Import(ctx, id, ecdsaPrivKeyb64, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.T().Run("should get a key pair successfully", func(t *testing.T) {
		keyRetrieved, err := s.store.Get(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, keyRetrieved.ID)
		assert.Equal(t, "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", key.PublicKey)
		assert.Equal(t, tags, keyRetrieved.Tags)
		assert.Equal(t, entities.Secp256k1, keyRetrieved.Algo.EllipticCurve)
		assert.Equal(t, entities.Ecdsa, keyRetrieved.Algo.Type)
		assert.Equal(t, "1", keyRetrieved.Metadata.Version)
		assert.NotNil(t, keyRetrieved.Metadata.CreatedAt)
		assert.NotNil(t, keyRetrieved.Metadata.UpdatedAt)
		assert.True(t, keyRetrieved.Metadata.DeletedAt.IsZero())
		assert.True(t, keyRetrieved.Metadata.DestroyedAt.IsZero())
		assert.True(t, keyRetrieved.Metadata.ExpireAt.IsZero())
		assert.False(t, keyRetrieved.Metadata.Disabled)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		keyRetrieved, err := s.store.Get(ctx, "invalidID")

		require.Nil(t, keyRetrieved)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *hashicorpKeyTestSuite) TestList() {
	ctx := s.env.ctx
	id1 := "my-key-list1"
	id2 := "my-key-list2"
	tags := testutils.FakeTags()

	_, err := s.store.Import(ctx, id1, ecdsaPrivKeyb64, &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	_, err = s.store.Import(ctx, id2, eddsaPrivKeyb64, &entities.Algorithm{
		Type:          entities.Eddsa,
		EllipticCurve: entities.Bn254,
	}, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.T().Run("should list all key pairs", func(t *testing.T) {
		keys, err := s.store.List(ctx)
		require.NoError(t, err)

		assert.Contains(t, keys, id1)
		assert.Contains(t, keys, id2)
	})
}

func (s *hashicorpKeyTestSuite) TestSign() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	payload := hexutil.Encode([]byte("my data to sign"))

	s.T().Run("should sign a message successfully: ECDSA/Secp256k1", func(t *testing.T) {
		id := "my-key-sign-ecdsa"

		_, err := s.store.Import(ctx, id, ecdsaPrivKeyb64, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		signature, err := s.store.Sign(ctx, id, payload)
		require.NoError(t, err)

		assert.Equal(t, "UWzxLZM7kztXXJGhWlkK0LeuYObJH7EOnMjv48qs6GB5rj7iEghkh3FfQyVCheWDTIHfdzBOst3eDRt0BGpaTg==", signature)
	})

	s.T().Run("should sign a message successfully: EDDSA/BN254", func(t *testing.T) {
		id := "my-key-sign-eddsa"

		_, err := s.store.Import(ctx, id, eddsaPrivKeyb64, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		signature, err := s.store.Sign(ctx, id, payload)
		require.NoError(t, err)

		assert.Equal(t, "RypSRagTLbR6tlOXu-REakfQRqRufPRCT8FxpZXuXZMDgwa5qYd5FAl1pRlLmQ_-alt1Ba4dKojknaVyHvCDeQ==", signature)
	})

	s.T().Run("should fail and parse the error code correctly", func(t *testing.T) {
		id := "my-key"

		key, err := s.store.Sign(ctx, id, "")

		require.Empty(t, key)
		assert.True(t, errors.IsInvalidFormatError(err))
	})
}
