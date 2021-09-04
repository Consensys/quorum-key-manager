package local

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mocksecrets "github.com/consensys/quorum-key-manager/src/stores/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

const (
	id             = "my-key"
	publicKeyECDSA = "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"
	publicKeyEDDSA = "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191"
	privKeyECDSA   = "0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
	privKeyEDDSA   = "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7"
)

var expectedErr = errors.DependencyFailureError("error")

type localKeyStoreTestSuite struct {
	suite.Suite
	keyStore        stores.KeyStore
	mockSecretStore *mocksecrets.MockSecretStore
}

func TestLocalKeyStore(t *testing.T) {
	s := new(localKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *localKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockSecretStore = mocksecrets.NewMockSecretStore(ctrl)

	s.keyStore = New(s.mockSecretStore, testutils2.NewMockLogger(ctrl))
}

func (s *localKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attr := testutils.FakeAttributes()

	s.Run("should create an ECDSA/Secp256k1 key successfully", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(testutils.FakeSecret(), nil)

		key, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, attr)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), attr.Tags, key.Tags)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
	})

	s.Run("should create an EDDSA/BN254 key successfully", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(testutils.FakeSecret(), nil)

		key, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, attr)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.NotEmpty(s.T(), key.PublicKey)
		assert.Equal(s.T(), attr.Tags, key.Tags)
		assert.Equal(s.T(), entities.Eddsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Bn254, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
	})

	s.Run("should fail with same error if Set fails", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(nil, expectedErr)

		_, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, attr)

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestImport() {
	ctx := context.Background()
	attr := testutils.FakeAttributes()

	s.Run("should create an ECDSA/Secp256k1 key successfully", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(testutils.FakeSecret(), nil)

		key, err := s.keyStore.Import(ctx, id, hexutil.MustDecode(privKeyECDSA), &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		}, attr)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), publicKeyECDSA, hexutil.Encode(key.PublicKey))
		assert.Equal(s.T(), attr.Tags, key.Tags)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
	})

	s.Run("should create an EDDSA/BN254 key successfully", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(testutils.FakeSecret(), nil)

		key, err := s.keyStore.Import(ctx, id, hexutil.MustDecode(privKeyEDDSA), &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, attr)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), publicKeyEDDSA, hexutil.Encode(key.PublicKey))
		assert.Equal(s.T(), attr.Tags, key.Tags)
		assert.Equal(s.T(), entities.Eddsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Bn254, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.UpdatedAt)
	})

	s.Run("should fail with InvalidParameter if algo is undefined", func() {
		_, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          "wrongType",
			EllipticCurve: entities.Bn254,
		}, attr)

		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})

	s.Run("should fail with same error if Set fails", func() {
		s.mockSecretStore.EXPECT().Set(ctx, id, gomock.Any(), &entities.Attributes{}).Return(nil, expectedErr)

		_, err := s.keyStore.Create(ctx, id, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		}, attr)

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestSign() {
	ctx := context.Background()

	s.Run("should sign with an ECDSA/Secp256k1 key successfully", func() {
		payload := crypto.Keccak256([]byte("my data"))
		secret := testutils.FakeSecret()
		secret.Value = base64.StdEncoding.EncodeToString(hexutil.MustDecode(privKeyECDSA))

		s.mockSecretStore.EXPECT().Get(ctx, id, "").Return(secret, nil)

		signature, err := s.keyStore.Sign(ctx, id, payload, &entities.Algorithm{
			Type:          entities.Ecdsa,
			EllipticCurve: entities.Secp256k1,
		})
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), "xUBOm7wht727RjpUY+KqK/NpCIOkzxX9H+dSBIWOITccTl/i5DyFvrcO3EIZTLV1gLVfCL+AOkY2pGWnIxygtQ==", base64.StdEncoding.EncodeToString(signature))
	})

	s.Run("should create an EDDSA/BN254 key successfully", func() {
		payload := []byte("my data")
		secret := testutils.FakeSecret()
		secret.Value = base64.StdEncoding.EncodeToString(hexutil.MustDecode(privKeyEDDSA))

		s.mockSecretStore.EXPECT().Get(ctx, id, "").Return(secret, nil)

		signature, err := s.keyStore.Sign(ctx, id, payload, &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		})
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), "YSmChRZfnuMYdhF8MJI46uy3W1aO6P2QV4Ed//kTCIQCoapFsBmln7bUviIjy12/diRYK/ZS70gjGYa+qHK11w==", base64.StdEncoding.EncodeToString(signature))
	})

	s.Run("should fail with InvalidParameter if algo is undefined", func() {
		payload := []byte("my data")
		secret := testutils.FakeSecret()
		secret.Value = base64.StdEncoding.EncodeToString(hexutil.MustDecode(privKeyECDSA))

		s.mockSecretStore.EXPECT().Get(ctx, id, "").Return(secret, nil)

		_, err := s.keyStore.Sign(ctx, id, payload, &entities.Algorithm{
			Type:          "wrongType",
			EllipticCurve: entities.Bn254,
		})

		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})

	s.Run("should fail with same error if Get fails", func() {
		s.mockSecretStore.EXPECT().Get(ctx, id, "").Return(nil, expectedErr)

		_, err := s.keyStore.Sign(ctx, id, []byte("my data"), &entities.Algorithm{
			Type:          entities.Eddsa,
			EllipticCurve: entities.Bn254,
		})

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestUpdate() {
	ctx := context.Background()

	s.Run("should return NotSupportedError", func() {
		_, err := s.keyStore.Update(ctx, id, testutils.FakeAttributes())
		assert.Equal(s.T(), errors.ErrNotSupported, err)
	})
}

func (s *localKeyStoreTestSuite) TestDelete() {
	ctx := context.Background()

	s.Run("should delete a key successfully", func() {
		s.mockSecretStore.EXPECT().Delete(ctx, id).Return(nil)

		err := s.keyStore.Delete(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Delete Secret fails", func() {
		s.mockSecretStore.EXPECT().Delete(ctx, id).Return(expectedErr)

		err := s.keyStore.Delete(ctx, id)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestUndelete() {
	ctx := context.Background()

	s.Run("should delete a key successfully", func() {
		s.mockSecretStore.EXPECT().Restore(ctx, id).Return(nil)

		err := s.keyStore.Restore(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Delete Secret fails", func() {
		s.mockSecretStore.EXPECT().Restore(ctx, id).Return(expectedErr)

		err := s.keyStore.Restore(ctx, id)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestDestroy() {
	ctx := context.Background()

	s.Run("should delete a key successfully", func() {
		s.mockSecretStore.EXPECT().Destroy(ctx, id).Return(nil)

		err := s.keyStore.Destroy(ctx, id)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Delete Secret fails", func() {
		s.mockSecretStore.EXPECT().Destroy(ctx, id).Return(expectedErr)

		err := s.keyStore.Destroy(ctx, id)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *localKeyStoreTestSuite) TestEncrypt() {
	ctx := context.Background()

	s.Run("should return NotImplementedError", func() {
		_, err := s.keyStore.Encrypt(ctx, id, []byte(""))
		assert.Equal(s.T(), errors.ErrNotImplemented, err)
	})
}

func (s *localKeyStoreTestSuite) TestDecrypt() {
	ctx := context.Background()

	s.Run("should return NotImplementedError", func() {
		_, err := s.keyStore.Decrypt(ctx, id, []byte(""))
		assert.Equal(s.T(), errors.ErrNotImplemented, err)
	})
}
