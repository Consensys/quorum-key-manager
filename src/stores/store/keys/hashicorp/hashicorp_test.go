package hashicorp

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/mocks"
	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/golang/mock/gomock"
	hashicorp "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id        = "my-key"
	publicKey = "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI="
)

var expectedErr = errors.HashicorpVaultError("error")

type hashicorpKeyStoreTestSuite struct {
	suite.Suite
	mockVault *mocks.MockPluginClient
	keyStore  stores.KeyStore
}

func TestHashicorpKeyStore(t *testing.T) {
	s := new(hashicorpKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockVault = mocks.NewMockPluginClient(ctrl)

	s.keyStore = New(s.mockVault, testutils2.NewMockLogger(ctrl))
}

func (s *hashicorpKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	expectedData := map[string]interface{}{
		idLabel:        id,
		curveLabel:     algorithm.EllipticCurve,
		algorithmLabel: algorithm.Type,
		tagsLabel:      attributes.Tags,
	}
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			idLabel:        id,
			publicKeyLabel: publicKey,
			curveLabel:     string(entities.Secp256k1),
			algorithmLabel: string(entities.Ecdsa),
			tagsLabel: map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			createdAtLabel: time.Now().Format(time.RFC3339),
			updatedAtLabel: time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should create a new key successfully", func() {
		s.mockVault.EXPECT().CreateKey(expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), attributes.Tags, key.Tags)
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail with NotSupported error", func() {
		s.mockVault.EXPECT().CreateKey(expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Curve25519,
		}, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsNotSupportedError(err))
	})

	s.Run("should fail with same error if CreateKey fails", func() {
		s.mockVault.EXPECT().CreateKey(expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestImport() {
	ctx := context.Background()
	privKey := "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="
	privKeyB, _ := base64.URLEncoding.DecodeString(privKey)
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	expectedData := map[string]interface{}{
		idLabel:         id,
		curveLabel:      algorithm.EllipticCurve,
		algorithmLabel:  algorithm.Type,
		tagsLabel:       attributes.Tags,
		privateKeyLabel: privKey,
	}
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			idLabel:        id,
			publicKeyLabel: publicKey,
			curveLabel:     string(entities.Secp256k1),
			algorithmLabel: string(entities.Ecdsa),
			tagsLabel: map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			createdAtLabel: time.Now().Format(time.RFC3339),
			updatedAtLabel: time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should import a new key successfully", func() {
		s.mockVault.EXPECT().ImportKey(expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), attributes.Tags, key.Tags)
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail with NotSupported error", func() {
		s.mockVault.EXPECT().CreateKey(expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Import(ctx, id, privKeyB, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Curve25519,
		}, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsNotSupportedError(err))
	})

	s.Run("should fail with same error if ImportKey fails", func() {
		s.mockVault.EXPECT().ImportKey(expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			idLabel:        id,
			publicKeyLabel: publicKey,
			curveLabel:     string(entities.Secp256k1),
			algorithmLabel: string(entities.Ecdsa),
			tagsLabel: map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			createdAtLabel: time.Now().Format(time.RFC3339),
			updatedAtLabel: time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should get a key successfully without version", func() {
		s.mockVault.EXPECT().GetKey(id).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Get(ctx, id)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), attributes.Tags, key.Tags)
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail with same error if GetKey fails", func() {
		s.mockVault.EXPECT().GetKey(id).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedIds := []interface{}{"my-key1", "my-key2"}

	s.Run("should list all secret ids successfully", func() {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": expectedIds,
			},
		}

		s.mockVault.EXPECT().ListKeys().Return(hashicorpSecret, nil)

		ids, err := s.keyStore.List(ctx, 0, 0)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{"my-key1", "my-key2"}, ids)
	})

	s.Run("should fail with same error if List fails", func() {
		s.mockVault.EXPECT().ListKeys().Return(nil, expectedErr)

		key, err := s.keyStore.List(ctx, 0, 0)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	data := []byte("my data")
	expectedSignature := base64.URLEncoding.EncodeToString([]byte("mySignature"))
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			signatureLabel: expectedSignature,
		},
	}

	s.Run("should sign a payload successfully", func() {
		s.mockVault.EXPECT().Sign(id, data).Return(hashicorpSecret, nil)

		signature, err := s.keyStore.Sign(ctx, id, data, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedSignature, base64.URLEncoding.EncodeToString(signature))
	})

	s.Run("should fail with NotSupported error", func() {
		signature, err := s.keyStore.Sign(ctx, id, data, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Curve25519,
		})

		assert.Nil(s.T(), signature)
		assert.True(s.T(), errors.IsNotSupportedError(err))
	})

	s.Run("should fail with same error if Sign fails", func() {
		s.mockVault.EXPECT().Sign(id, data).Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, data, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.Empty(s.T(), signature)
		require.NotNil(s.T(), err)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}
