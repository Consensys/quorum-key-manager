// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
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
	s.T().Run("should create a new key successfully: Secp256k1/ECDSA", func(t *testing.T) {
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
		require.NoError(t, err)

		assert.NotEmpty(t, key.PublicKey)
		assert.Equal(t, request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(t, request.Curve, key.Curve)
		assert.Equal(t, request.ID, key.ID)
		assert.Equal(t, request.Tags, key.Tags)
		assert.Equal(t, "1", key.Version)
		assert.False(t, key.Disabled)
		assert.NotEmpty(t, key.CreatedAt)
		assert.NotEmpty(t, key.UpdatedAt)
		assert.True(t, key.ExpireAt.IsZero())
		assert.True(t, key.DeletedAt.IsZero())
		assert.True(t, key.DestroyedAt.IsZero())
	})

	s.T().Run("should create a new key successfully: BN254/EDDSA", func(t *testing.T) {
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
		require.NoError(t, err)

		assert.NotEmpty(t, key.PublicKey)
		assert.Equal(t, request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(t, request.Curve, key.Curve)
		assert.Equal(t, request.ID, key.ID)
		assert.Equal(t, request.Tags, key.Tags)
		assert.Equal(t, "1", key.Version)
		assert.False(t, key.Disabled)
		assert.NotEmpty(t, key.CreatedAt)
		assert.NotEmpty(t, key.UpdatedAt)
		assert.True(t, key.ExpireAt.IsZero())
		assert.True(t, key.DeletedAt.IsZero())
		assert.True(t, key.DestroyedAt.IsZero())
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
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
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
		assert.Equal(t, " key store inexistentStoreName was not found", httpError.Message)
	})

	s.T().Run("should fail with bad request if curve is not supported", func(t *testing.T) {
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
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})

	s.T().Run("should fail with bad request if signing algorithm is not supported", func(t *testing.T) {
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
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestImport() {
	s.T().Run("should create a new key successfully: Secp256k1/ECDSA", func(t *testing.T) {
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
		require.NoError(t, err)

		assert.Equal(t, "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", key.PublicKey)
		assert.Equal(t, request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(t, request.Curve, key.Curve)
		assert.Equal(t, request.ID, key.ID)
		assert.Equal(t, request.Tags, key.Tags)
		assert.Equal(t, "1", key.Version)
		assert.False(t, key.Disabled)
		assert.NotEmpty(t, key.CreatedAt)
		assert.NotEmpty(t, key.UpdatedAt)
		assert.True(t, key.ExpireAt.IsZero())
		assert.True(t, key.DeletedAt.IsZero())
		assert.True(t, key.DestroyedAt.IsZero())
	})

	s.T().Run("should create a new key successfully: BN254/EDDSA", func(t *testing.T) {
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
		require.NoError(t, err)

		assert.Equal(t, "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", key.PublicKey)
		assert.Equal(t, request.SigningAlgorithm, key.SigningAlgorithm)
		assert.Equal(t, request.Curve, key.Curve)
		assert.Equal(t, request.ID, key.ID)
		assert.Equal(t, request.Tags, key.Tags)
		assert.Equal(t, "1", key.Version)
		assert.False(t, key.Disabled)
		assert.NotEmpty(t, key.CreatedAt)
		assert.NotEmpty(t, key.UpdatedAt)
		assert.True(t, key.ExpireAt.IsZero())
		assert.True(t, key.DeletedAt.IsZero())
		assert.True(t, key.DestroyedAt.IsZero())
	})

	s.T().Run("should fail with bad request if curve is not supported", func(t *testing.T) {
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
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})

	s.T().Run("should fail with bad request if signing algorithm is not supported", func(t *testing.T) {
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
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
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

	s.T().Run("should get a key successfully", func(t *testing.T) {
		keyRetrieved, err := s.keyManagerClient.GetKey(s.ctx, s.cfg.HashicorpKeyStore, key.ID)
		require.NoError(t, err)

		assert.Equal(t, "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", keyRetrieved.PublicKey)
		assert.Equal(t, request.ID, keyRetrieved.ID)
		assert.Equal(t, request.Tags, keyRetrieved.Tags)
		assert.Equal(t, "1", keyRetrieved.Version)
		assert.False(t, keyRetrieved.Disabled)
		assert.NotEmpty(t, keyRetrieved.CreatedAt)
		assert.NotEmpty(t, keyRetrieved.UpdatedAt)
		assert.True(t, keyRetrieved.ExpireAt.IsZero())
		assert.True(t, keyRetrieved.DeletedAt.IsZero())
		assert.True(t, keyRetrieved.DestroyedAt.IsZero())
	})

	// TODO: Add test to check that a specific version can be retrieved when versioning is implemented in Hashicorp

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		secret, err := s.keyManagerClient.GetSecret(s.ctx, s.cfg.HashicorpKeyStore, "inexistentID", key.Version)
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
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

	s.T().Run("should get all key ids successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListKeys(s.ctx, s.cfg.HashicorpKeyStore)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(ids), 1)
		assert.Contains(t, ids, key.ID)
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(s.ctx, "inexistentStoreName")
		require.Empty(t, ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
		assert.Equal(t, " secret store inexistentStoreName was not found", httpError.Message)
	})
}

func (s *keysTestSuite) TestSign() {
	s.T().Run("should sign a new payload successfully: Secp256k1/ECDSA", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-ecdsa",
			Curve:            "secp256k1",
			PrivateKey:       ecdsaPrivKey,
			SigningAlgorithm: "ecdsa",
		}

		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: hexutil.Encode([]byte("my data to sign")),
		}
		signature, err := s.keyManagerClient.Sign(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.NoError(t, err)

		assert.Equal(t, "UWzxLZM7kztXXJGhWlkK0LeuYObJH7EOnMjv48qs6GB5rj7iEghkh3FfQyVCheWDTIHfdzBOst3eDRt0BGpaTg==", signature)

	})

	s.T().Run("should sign a new payload successfully: BN254/EDDSA", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
		}
		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: hexutil.Encode([]byte("my data to sign")),
		}
		signature, err := s.keyManagerClient.Sign(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.NoError(t, err)

		assert.Equal(t, "RypSRagTLbR6tlOXu-REakfQRqRufPRCT8FxpZXuXZMDgwa5qYd5FAl1pRlLmQ_-alt1Ba4dKojknaVyHvCDeQ==", signature)
	})

	s.T().Run("should fail if payload is not base64 string", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       eddsaPrivKey,
		}
		key, err := s.keyManagerClient.ImportKey(s.ctx, s.cfg.HashicorpKeyStore, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: "my data to sign not in base64 format",
		}
		signature, err := s.keyManagerClient.Sign(s.ctx, s.cfg.HashicorpKeyStore, key.ID, requestSign)
		require.Empty(t, signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})
}
