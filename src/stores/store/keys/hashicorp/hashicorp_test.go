package hashicorp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/mocks"
	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
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
	mockVault  *mocks.MockVaultClient
	mountPoint string
	keyStore   stores.KeyStore
}

func TestHashicorpKeyStore(t *testing.T) {
	s := new(hashicorpKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *hashicorpKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mountPoint = "hashicorp-plugin"
	s.mockVault = mocks.NewMockVaultClient(ctrl)

	s.keyStore = New(s.mockVault, s.mountPoint, testutils2.NewMockLogger(ctrl))
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
			"id":         id,
			"public_key": publicKey,
			"curve":      string(entities.Secp256k1),
			"algorithm":  string(entities.Ecdsa),
			"tags": map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			"version":    json.Number("1"),
			"created_at": time.Now().Format(time.RFC3339),
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should create a new key successfully", func() {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

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

	s.Run("should fail with same error if write fails", func() {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestImport() {
	ctx := context.Background()
	privKey := "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="
	privKeyB, _ := base64.URLEncoding.DecodeString(privKey)
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
			"id":         id,
			"public_key": publicKey,
			"curve":      string(entities.Secp256k1),
			"algorithm":  string(entities.Ecdsa),
			"tags": map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			"version":    json.Number("1"),
			"created_at": time.Now().Format(time.RFC3339),
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should import a new key successfully", func() {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

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

	s.Run("should fail with same error if write fails", func() {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys/" + id
	attributes := testutils.FakeAttributes()
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			"id":         id,
			"public_key": publicKey,
			"curve":      string(entities.Secp256k1),
			"algorithm":  string(entities.Ecdsa),
			"tags": map[string]interface{}{
				"tag1": "tagValue1",
				"tag2": "tagValue2",
			},
			"version":    json.Number("1"),
			"created_at": time.Now().Format(time.RFC3339),
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}

	s.Run("should get a key successfully without version", func() {
		s.mockVault.EXPECT().Read(expectedPath, nil).Return(hashicorpSecret, nil)

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

	s.Run("should fail with same error if read fails", func() {
		s.mockVault.EXPECT().Read(expectedPath, nil).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys"
	expectedIds := []interface{}{"my-key1", "my-key2"}

	s.Run("should list all secret ids successfully", func() {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": expectedIds,
			},
		}

		s.mockVault.EXPECT().List(expectedPath).Return(hashicorpSecret, nil)

		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{"my-key1", "my-key2"}, ids)
	})

	s.Run("should fail with same error if read fails", func() {
		s.mockVault.EXPECT().List(expectedPath).Return(nil, expectedErr)

		key, err := s.keyStore.List(ctx)

		assert.Nil(s.T(), key)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}

func (s *hashicorpKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys/" + id + "/sign"
	data := []byte("my data")
	expectedData := base64.URLEncoding.EncodeToString(data)
	expectedSignature := base64.URLEncoding.EncodeToString([]byte("mySignature"))
	hashicorpSecret := &hashicorp.Secret{
		Data: map[string]interface{}{
			signatureLabel: expectedSignature,
		},
	}

	s.Run("should sign a payload successfully", func() {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(hashicorpSecret, nil)

		signature, err := s.keyStore.Sign(ctx, id, data, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedSignature, base64.URLEncoding.EncodeToString(signature))
	})

	s.Run("should fail with same error if write fails", func() {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, data, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})

		assert.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsHashicorpVaultError(err))
	})
}
