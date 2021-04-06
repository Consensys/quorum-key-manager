package hashicorp

import (
	"context"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/golang/mock/gomock"
	hashicorp "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id        = "my-key"
	publicKey = "0x0433d7f005495fb6c0a34e22336dc3adcf4064553d5e194f77126bcac6da19491e0bab2772115cd284605d3bba94b69dc8c7a215021b58bcc87a70c9a440a3ff83"
)

type hashicorpKeyStoreTestSuite struct {
	suite.Suite
	mockVault  *mocks.MockVaultClient
	mountPoint string
	keyStore   keys.Store
}

func TestHashicorpSecretStore(t *testing.T) {
	s := new(hashicorpKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mountPoint = "hashicorp-plugin"
	s.mockVault = mocks.NewMockVaultClient(ctrl)

	s.keyStore = New(s.mockVault, s.mountPoint)
}

func (s *hashicorpKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys"
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
			"id":        id,
			"publicKey": publicKey,
			"curve":     entities.Secp256k1,
			"algorithm": entities.Ecdsa,
			"tags":      testutils.FakeTags(),
		},
	}

	s.T().Run("should create a new key successfully", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, key.PublicKey)
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, 1, key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with NotFound error if write fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidFormat error if write fails with 400", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusBadRequest,
		})

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with InvalidParameter error if write fails with 422", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
		})

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with AlreadyExists error if write fails with 409", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusConflict,
		})

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should fail with HashicorpVaultConnection error if write fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestImport() {
	ctx := context.Background()
	privKey := "0b0232595b77568d99364bede133839ccbcb40775967a7eacd15d355c96288b5"
	expectedPath := s.mountPoint + "/keys/import"
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
			"id":        id,
			"publicKey": publicKey,
			"curve":     entities.Secp256k1,
			"algorithm": entities.Ecdsa,
			"tags":      testutils.FakeTags(),
		},
	}

	s.T().Run("should import a new key successfully", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, key.PublicKey)
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, 1, key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with NotFound error if write fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidFormat error if write fails with 400", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusBadRequest,
		})

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with InvalidParameter error if write fails with 422", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
		})

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with AlreadyExists error if write fails with 409", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusConflict,
		})

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should fail with HashicorpVaultConnection error if write fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		key, err := s.keyStore.Import(ctx, id, privKey, algorithm, attributes)

		assert.Nil(t, key)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys/" + id
	attributes := testutils.FakeAttributes()
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			"id":        id,
			"publicKey": publicKey,
			"curve":     entities.Secp256k1,
			"algorithm": entities.Ecdsa,
			"tags":      testutils.FakeTags(),
		},
	}

	s.T().Run("should get a key successfully without version", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPath).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Get(ctx, id, "")

		assert.NoError(t, err)
		assert.Equal(t, publicKey, key.PublicKey)
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, 1, key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with NotFound error if read fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPath).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})

		key, err := s.keyStore.Get(ctx, id, "")

		assert.Nil(t, key)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with HashicorpVaultConnection error if read fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPath).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		key, err := s.keyStore.Get(ctx, id, "")

		assert.Nil(t, key)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys"
	expectedIds := []string{"my-key1", "my-key2"}

	s.T().Run("should list all secret ids successfully", func(t *testing.T) {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": expectedIds,
			},
		}

		s.mockVault.EXPECT().List(expectedPath).Return(hashicorpSecret, nil)

		ids, err := s.keyStore.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedIds, ids)
	})

	s.T().Run("should fail with HashicorpVaultConnection error if list fails", func(t *testing.T) {
		s.mockVault.EXPECT().List(expectedPath).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		key, err := s.keyStore.List(ctx)

		assert.Nil(t, key)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys/" + id + "/sign"
	expectedData := "my data"
	expectedSignature := "0x8b9679a75861e72fa6968dd5add3bf96e2747f0f124a2e728980f91e1958367e19c2486a40fdc65861824f247603bc18255fa497ca0b8b0a394aa7a6740fdc4601"
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			signatureLabel: expectedSignature,
		},
	}

	s.T().Run("should refresh a secret without expiration date", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(hashicorpSecret, nil)

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, signature)
	})

	s.T().Run("should fail with NotFound error if write fails with 404", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusNotFound,
		})

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.Empty(t, signature)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidFormat error if write fails with 400", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusBadRequest,
		})

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with InvalidParameter error if write fails with 422", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
		})

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with AlreadyExists error if write fails with 409", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusConflict,
		})

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.Empty(t, signature)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should fail with HashicorpVaultConnection error if write fails with 500", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, &hashicorp.ResponseError{
			StatusCode: http.StatusInternalServerError,
		})

		signature, err := s.keyStore.Sign(ctx, id, expectedData, "")

		assert.Empty(t, signature)
		assert.True(t, errors.IsHashicorpVaultConnectionError(err))
	})
}
