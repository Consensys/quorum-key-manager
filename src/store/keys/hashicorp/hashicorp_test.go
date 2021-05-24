package hashicorp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
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
	publicKey = "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI="
)

type hashicorpKeyStoreTestSuite struct {
	suite.Suite
	mockVault  *mocks.MockVaultClient
	mountPoint string
	keyStore   keys.Store
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

	logger := log.DefaultLogger()
	s.keyStore = New(s.mockVault, s.mountPoint, logger)
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

	s.T().Run("should create a new key successfully", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, "1", key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
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

	s.T().Run("should import a new key successfully", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, "1", key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockVault.EXPECT().Write(expectedPath, expectedData).Return(nil, expectedErr)

		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
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

	s.T().Run("should get a key successfully without version", func(t *testing.T) {
		s.mockVault.EXPECT().Read(expectedPath, nil).Return(hashicorpSecret, nil)

		key, err := s.keyStore.Get(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, base64.URLEncoding.EncodeToString(key.PublicKey))
		assert.Equal(t, id, key.ID)
		assert.Equal(t, entities.Ecdsa, key.Algo.Type)
		assert.Equal(t, entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(t, key.Metadata.Disabled)
		assert.Equal(t, "1", key.Metadata.Version)
		assert.Equal(t, attributes.Tags, key.Tags)
		assert.True(t, key.Metadata.ExpireAt.IsZero())
		assert.True(t, key.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if read fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockVault.EXPECT().Read(expectedPath, nil).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
	})
}

func (s *hashicorpKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedPath := s.mountPoint + "/keys"
	expectedIds := []interface{}{"my-key1", "my-key2"}

	s.T().Run("should list all secret ids successfully", func(t *testing.T) {
		hashicorpSecret := &hashicorp.Secret{
			Data: map[string]interface{}{
				"keys": expectedIds,
			},
		}

		s.mockVault.EXPECT().List(expectedPath).Return(hashicorpSecret, nil)

		ids, err := s.keyStore.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"my-key1", "my-key2"}, ids)
	})

	s.T().Run("should fail with same error if read fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockVault.EXPECT().List(expectedPath).Return(nil, expectedErr)

		key, err := s.keyStore.List(ctx)

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
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

	s.T().Run("should sign a payload successfully", func(t *testing.T) {
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(hashicorpSecret, nil)

		signature, err := s.keyStore.Sign(ctx, id, data)

		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, base64.URLEncoding.EncodeToString(signature))
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockVault.EXPECT().Write(expectedPath, map[string]interface{}{
			dataLabel: expectedData,
		}).Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, data)

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})
}
