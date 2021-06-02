// +build e2e

package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	ecdsaPrivKey = "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="
	eddsaPrivKey = "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZGGbioek5qYuzJzTNZpTHrVjjFk7iFe3FYwfpxZyNPxtIaFB5gb9VP9IcHZewwNZly821re7RkmB8pGdjywygPH"
)

type keysTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
}

func (s *keysTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})
}

func (s *keysTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerKeys(t *testing.T) {
	s := new(keysTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *keysTestSuite) TestCreate() {
	s.Run("should create a new key successfully: Secp256k1/ECDSA", func() {
		request := &types.CreateKeyRequest{
			ID:               "my-key-ecdsa",
			Curve:            "secp256k1",
			SigningAlgorithm: "ecdsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.CreateKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), request.ID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
		assert.True(s.T(), key.ExpireAt.IsZero())
		assert.True(s.T(), key.DeletedAt.IsZero())
		assert.True(s.T(), key.DestroyedAt.IsZero())
	})

	s.Run("should create a new key successfully: BN254/EDDSA", func() {
		request := &types.CreateKeyRequest{
			ID:               "my-key-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.CreateKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), request.ID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
		assert.True(s.T(), key.ExpireAt.IsZero())
		assert.True(s.T(), key.DeletedAt.IsZero())
		assert.True(s.T(), key.DestroyedAt.IsZero())
	})

	s.Run("should parse errors successfully", func() {
		request := &types.CreateKeyRequest{
			ID:               "my-key",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.CreateKey(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.Run("should fail with bad request if curve is not supported", func() {
		request := &types.CreateKeyRequest{
			ID:               "my-key",
			Curve:            "invalidCurve",
			SigningAlgorithm: "eddsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.CreateKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})

	s.Run("should fail with bad request if signing algorithm is not supported", func() {
		request := &types.CreateKeyRequest{
			ID:               "my-key",
			Curve:            "secp256k1",
			SigningAlgorithm: "invalidSigningAlgorithm",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.CreateKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestImport() {
	s.Run("should create a new key successfully: Secp256k1/ECDSA", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import-ecdsa",
			Curve:            "secp256k1",
			PrivateKey:       ecdsaPrivKey,
			SigningAlgorithm: "ecdsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), request.ID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
		assert.True(s.T(), key.ExpireAt.IsZero())
		assert.True(s.T(), key.DeletedAt.IsZero())
		assert.True(s.T(), key.DestroyedAt.IsZero())
	})

	s.Run("should create a new key successfully: BN254/EDDSA", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", key.PublicKey)
		assert.Equal(s.T(), request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(s.T(), request.Curve, key.Curve)
		assert.Equal(s.T(), request.ID, key.ID)
		assert.Equal(s.T(), request.Tags, key.Tags)
		assert.False(s.T(), key.Disabled)
		assert.NotEmpty(s.T(), key.CreatedAt)
		assert.NotEmpty(s.T(), key.UpdatedAt)
		assert.True(s.T(), key.ExpireAt.IsZero())
		assert.True(s.T(), key.DeletedAt.IsZero())
		assert.True(s.T(), key.DestroyedAt.IsZero())
	})

	s.Run("should fail with bad request if curve is not supported", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import",
			Curve:            "invalidCurve",
			SigningAlgorithm: "eddsa",
			PrivateKey:       ecdsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})

	s.Run("should fail with bad request if signing algorithm is not supported", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import",
			Curve:            "secp256k1",
			SigningAlgorithm: "invalidSigningAlgorithm",
			PrivateKey:       ecdsaPrivKey,
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestGet() {
	request := &types.ImportKeyRequest{
		ID:               "my-key-get",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
	require.NoError(s.T(), err)

	s.Run("should get a key successfully", func() {
		keyRetrieved, err := s.keyManagerClient.GetKey(s.ctx, s.cfg.HashicorpKeyStore, key.ID)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", keyRetrieved.PublicKey)
		assert.Equal(s.T(), request.ID, keyRetrieved.ID)
		assert.Equal(s.T(), request.Tags, keyRetrieved.Tags)
		assert.False(s.T(), keyRetrieved.Disabled)
		assert.NotEmpty(s.T(), keyRetrieved.CreatedAt)
		assert.NotEmpty(s.T(), keyRetrieved.UpdatedAt)
		assert.True(s.T(), keyRetrieved.ExpireAt.IsZero())
		assert.True(s.T(), keyRetrieved.DeletedAt.IsZero())
		assert.True(s.T(), keyRetrieved.DestroyedAt.IsZero())
	})
}

func (s *keysTestSuite) TestList() {
	request := &types.ImportKeyRequest{
		ID:               "my-key-list",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       ecdsaPrivKey,
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
	require.NoError(s.T(), err)

	s.Run("should get all key ids successfully", func() {
		ids, err := s.keyManagerClient.ListKeys(s.ctx, s.cfg.HashicorpKeyStore)
		require.NoError(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(ids), 1)
		assert.Contains(s.T(), ids, key.ID)
	})

	s.Run("should parse errors successfully", func() {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, "inexistentStoreName")
		require.Empty(s.T(), ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestSign() {
	data := []byte("my data to sign")
	hashedPayload := base64.URLEncoding.EncodeToString(crypto.Keccak256(data))
	payload := base64.URLEncoding.EncodeToString(data)

	s.Run("should sign a new payload successfully: Secp256k1/ECDSA", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-ecdsa",
			Curve:            "secp256k1",
			PrivateKey:       ecdsaPrivKey,
			SigningAlgorithm: "ecdsa",
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		requestSign := &types.SignBase64PayloadRequest{
			Data: hashedPayload,
		}
		signature, err := s.keyManagerClient.SignKey(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "YzQeLIN0Sd43Nbb0QCsVSqChGNAuRaKzEfujnERAJd0523aZyz2KXK93KKh-d4ws3MxAhc8qNG43wYI97Fzi7Q==", signature)

	})

	s.Run("should sign a new payload successfully: BN254/EDDSA", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
		}
		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		requestSign := &types.SignBase64PayloadRequest{
			Data: payload,
		}
		signature, err := s.keyManagerClient.SignKey(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), "tdpR9JkX7lKSugSvYJX2icf6_uQnCAmXG9v_FG26vS0AcBqg6eVakZQNYwfic_Ec3LWqzSbXg54TBteQq6grdw==", signature)
	})

	s.Run("should fail if payload is not base64 string", func() {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
		}
		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(s.T(), err)

		requestSign := &types.SignBase64PayloadRequest{
			Data: "my data to sign not in base64 format",
		}
		signature, err := s.keyManagerClient.SignKey(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 400, httpError.StatusCode)
	})
}
