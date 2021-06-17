package akv

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/log/mock"

	akv "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/akv/mocks"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	id        = "my-key"
	publicKey = "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"
	privKey   = "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
)

var (
	base64PrivKey = "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw"
	base64PubKeyX = "VVIUmGpSH0NAnBxrI22xZ0My-qrxH8QqcEerB3gevm8"
	base64PubKeyY = "CXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI"
)

type akvKeyStoreTestSuite struct {
	suite.Suite
	mockVault *mocks.MockKeysClient
	keyStore  keys.Store
}

func TestHashicorpKeyStore(t *testing.T) {
	s := new(akvKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *akvKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockVault = mocks.NewMockKeysClient(ctrl)
	s.keyStore = New(s.mockVault, mock.NewMockLogger(ctrl))
}

func (s *akvKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	version := "1234"

	akvKeyID := fmt.Sprintf("keyvault.com/keys/%s/%s", id, version)
	akvKey := akv.KeyBundle{
		Attributes: &akv.KeyAttributes{
			Enabled: common.ToPtr(true).(*bool),
			Created: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
			Updated: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
		},
		Key: &akv.JSONWebKey{
			Kid: &akvKeyID,
			Crv: akv.P256K,
			Kty: akv.EC,
			X:   &base64PubKeyX,
			Y:   &base64PubKeyY,
		},
	}

	s.Run("should create a new key successfully", func() {
		s.mockVault.EXPECT().CreateKey(gomock.Any(), id, akv.EC, akv.P256K, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(akvKey, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, hexutil.Encode(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), version, key.Metadata.Version)
	})
}

func (s *akvKeyStoreTestSuite) TestImport() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	version := "1234"

	akvKeyID := fmt.Sprintf("keyvault.com/keys/%s/%s", id, version)
	akvKey := akv.KeyBundle{
		Attributes: &akv.KeyAttributes{
			Enabled: common.ToPtr(true).(*bool),
			Created: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
			Updated: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
		},
		Key: &akv.JSONWebKey{
			Kid: &akvKeyID,
			Crv: akv.P256K,
			Kty: akv.EC,
			X:   &base64PubKeyX,
			Y:   &base64PubKeyY,
		},
	}

	s.Run("should create a new key successfully", func() {
		s.mockVault.EXPECT().ImportKey(gomock.Any(), id, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, keyName string, k *akv.JSONWebKey, attr *akv.KeyAttributes, tags map[string]string) (akv.KeyBundle, error) {
				require.Equal(s.T(), k.Crv, akv.P256K)
				require.Equal(s.T(), k.Kty, akv.EC)
				require.Equal(s.T(), *k.D, base64PrivKey)
				require.Equal(s.T(), *k.X, base64PubKeyX)
				require.Equal(s.T(), *k.Y, base64PubKeyY)
				return akvKey, nil
			})

		privKeyB, _ := hex.DecodeString(privKey)
		key, err := s.keyStore.Import(ctx, id, privKeyB, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, hexutil.Encode(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), version, key.Metadata.Version)
	})
}

func (s *akvKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	version := "1234"

	akvKeyID := fmt.Sprintf("keyvault.com/keys/%s/%s", id, version)
	akvKey := akv.KeyBundle{
		Attributes: &akv.KeyAttributes{
			Enabled: common.ToPtr(true).(*bool),
			Created: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
			Updated: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
		},
		Tags: common.Tomapstrptr(attributes.Tags),
		Key: &akv.JSONWebKey{
			Kid: &akvKeyID,
			Crv: akv.P256K,
			Kty: akv.EC,
			X:   &base64PubKeyX,
			Y:   &base64PubKeyY,
		},
	}

	s.Run("should get a key successfully", func() {
		s.mockVault.EXPECT().GetKey(gomock.Any(), id, "").Return(akvKey, nil)

		key, err := s.keyStore.Get(ctx, id)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), publicKey, hexutil.Encode(key.PublicKey))
		assert.Equal(s.T(), id, key.ID)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), version, key.Metadata.Version)
		assert.Equal(s.T(), attributes.Tags, key.Tags)
		assert.True(s.T(), key.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), key.Metadata.DeletedAt.IsZero())
	})
}

func (s *akvKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	expectedIds := []interface{}{"my-key1", "my-key2"}
	kIds := []string{"myvault.com/keys/" + expectedIds[0].(string), "myvault.com/keys/" + expectedIds[1].(string)}

	s.Run("should list all secret ids successfully", func() {
		keyList := []akv.KeyItem{{Kid: &kIds[0]}, {Kid: &kIds[1]}}

		s.mockVault.EXPECT().GetKeys(gomock.Any(), gomock.Any()).Return(keyList, nil)

		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{"my-key1", "my-key2"}, ids)
	})
}

func (s *akvKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	version := "1234"
	payload := []byte("my data")
	attributes := testutils.FakeAttributes()
	expectedSignature := "0x8b9679a75861e72fa6968dd5add3bf96e2747f0f124a2e728980f91e1958367e19c2486a40fdc65861824f247603bc18255fa497ca0b8b0a394aa7a6740fdc4601"
	akvKeyID := fmt.Sprintf("keyvault.com/keys/%s/%s", id, version)

	akvKey := akv.KeyBundle{
		Attributes: &akv.KeyAttributes{
			Enabled: common.ToPtr(true).(*bool),
			Created: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
			Updated: common.ToPtr(date.NewUnixTimeFromNanoseconds(time.Now().UnixNano())).(*date.UnixTime),
		},
		Tags: common.Tomapstrptr(attributes.Tags),
		Key: &akv.JSONWebKey{
			Kid: &akvKeyID,
			Crv: akv.P256K,
			Kty: akv.EC,
			X:   &base64PubKeyX,
			Y:   &base64PubKeyY,
		},
	}

	bSig, _ := hexutil.Decode(expectedSignature)
	b64Sig := base64.RawURLEncoding.EncodeToString(bSig)
	b64Payload := base64.RawURLEncoding.EncodeToString(payload)

	s.Run("should sign payload successfully", func() {
		s.mockVault.EXPECT().GetKey(gomock.Any(), id, "").Return(akvKey, nil)
		s.mockVault.EXPECT().Sign(gomock.Any(), id, "", akv.ES256K, b64Payload).Return(b64Sig, nil)

		signature, err := s.keyStore.Sign(ctx, id, payload)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), hexutil.Encode(signature), expectedSignature)
	})
}
