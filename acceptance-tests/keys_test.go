// +build acceptance

package integrationtests

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/suite"
)

type keysTestSuite struct {
	suite.Suite
	env              *IntegrationEnvironment
	err              error
	keyManagerClient *client.HTTPClient
}

func (s *keysTestSuite) SetupSuite() {
	err := StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.T().Error(err)
		return
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *keysTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerKeys(t *testing.T) {
	s := new(keysTestSuite)

	var err error
	s.env, err = NewIntegrationEnvironment(context.Background())
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.env.Cancel()
	})
	defer sig.Close()

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.env.baseURL,
	})

	suite.Run(t, s)
}

func (s *keysTestSuite) TestCreate() {
	ctx := s.env.ctx

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

		key, err := s.keyManagerClient.CreateKey(ctx, KeyStoreName, request)
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

		key, err := s.keyManagerClient.CreateKey(ctx, KeyStoreName, request)
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

		key, err := s.keyManagerClient.CreateKey(ctx, "inexistentStoreName", request)
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

		key, err := s.keyManagerClient.CreateKey(ctx, KeyStoreName, request)
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

		key, err := s.keyManagerClient.CreateKey(ctx, KeyStoreName, request)
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestImport() {
	ctx := s.env.ctx

	s.T().Run("should create a new key successfully: Secp256k1/ECDSA", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import-ecdsa",
			Curve:            "secp256k1",
			PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
			SigningAlgorithm: "ecdsa",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.NoError(t, err)

		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", key.PublicKey)
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
			PrivateKey:       "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.NoError(t, err)

		assert.Equal(t, "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191", key.PublicKey)
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
			PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})

	s.T().Run("should fail with bad request if signing algorithm is not supported", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-import",
			Curve:            "secp256k1",
			SigningAlgorithm: "invalidSigningAlgorithm",
			PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
			Tags: map[string]string{
				"myTag0": "tag0",
				"myTag1": "tag1",
			},
		}

		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.Nil(t, key)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestGet() {
	ctx := s.env.ctx
	request := &types.ImportKeyRequest{
		ID:               "my-key-get",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
	require.NoError(s.T(), err)

	s.T().Run("should get a key successfully", func(t *testing.T) {
		keyRetrieved, err := s.keyManagerClient.GetKey(ctx, KeyStoreName, key.ID, "")
		require.NoError(t, err)

		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", keyRetrieved.PublicKey)
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
		secret, err := s.keyManagerClient.GetSecret(ctx, KeyStoreName, "inexistentID", key.Version)
		require.Nil(t, secret)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
	})
}

func (s *keysTestSuite) TestList() {
	ctx := s.env.ctx
	request := &types.ImportKeyRequest{
		ID:               "my-key-list",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		PrivateKey:       "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249",
		Tags: map[string]string{
			"myTag0": "tag0",
			"myTag1": "tag1",
		},
	}

	key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
	require.NoError(s.T(), err)

	s.T().Run("should get all key ids successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListKeys(ctx, KeyStoreName)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(ids), 1)
		assert.Contains(t, ids, key.ID)
	})

	s.T().Run("should parse errors successfully", func(t *testing.T) {
		ids, err := s.keyManagerClient.ListSecrets(ctx, "inexistentStoreName")
		require.Empty(t, ids)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 404, httpError.StatusCode)
		assert.Equal(t, " secret store inexistentStoreName was not found", httpError.Message)
	})
}

func (s *keysTestSuite) TestSign() {
	ctx := s.env.ctx

	s.T().Run("should sign a new payload successfully: Secp256k1/ECDSA", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-ecdsa",
			Curve:            "secp256k1",
			PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
			SigningAlgorithm: "ecdsa",
		}

		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: hexutil.Encode([]byte("my data to sign")),
		}
		signature, err := s.keyManagerClient.Sign(ctx, KeyStoreName, key.ID, requestSign)
		require.NoError(t, err)

		assert.Equal(t, "0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed01", signature)

	})

	s.T().Run("should sign a new payload successfully: BN254/EDDSA", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7",
		}
		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: hexutil.Encode([]byte("my data to sign")),
		}
		signature, err := s.keyManagerClient.Sign(ctx, KeyStoreName, key.ID, requestSign)
		require.NoError(t, err)

		assert.Equal(t, "0xb5da51f49917ee5292ba04af6095f689c7fafee4270809971bdbff146dbabd2d00701aa0e9e55a91940d6307e273f11cdcb5aacd26d7839e1306d790aba82b77", signature)
	})

	s.T().Run("should fail if payload is not hexadecimal string", func(t *testing.T) {
		request := &types.ImportKeyRequest{
			ID:               "my-key-sign-eddsa",
			Curve:            "bn254",
			SigningAlgorithm: "eddsa",
			PrivateKey:       "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7",
		}
		key, err := s.keyManagerClient.ImportKey(ctx, KeyStoreName, request)
		require.NoError(t, err)

		requestSign := &types.SignPayloadRequest{
			Data: "my data to sign not in hexadecimal format",
		}
		signature, err := s.keyManagerClient.Sign(ctx, KeyStoreName, key.ID, requestSign)
		require.Empty(t, signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(t, 400, httpError.StatusCode)
	})
}
