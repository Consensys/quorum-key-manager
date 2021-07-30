package aws

import (
	"context"
	"encoding/base64"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/mock"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/infra/aws/mocks"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id    = "my-key"
	keyID = "key-ID"
)

var expectedErr = errors.AWSError("error")

type awsKeyStoreTestSuite struct {
	suite.Suite
	mockKmsClient *mocks.MockKmsClient
	mockKeys      *mock.MockKeys
	keyStore      keys.Store
}

func TestAWSKeyStore(t *testing.T) {
	s := new(awsKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *awsKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockKeys = mock.NewMockKeys(ctrl)

	s.mockKmsClient = mocks.NewMockKmsClient(ctrl)
	s.keyStore = New(s.mockKmsClient, s.mockKeys, testutils.NewMockLogger(ctrl))
}

func (s *awsKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils2.FakeAttributes()
	algorithm := testutils2.FakeAlgorithm()

	retCreateKey := kms.CreateKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId: aws.String(keyID),
		},
	}
	retGetPubKey := fakeGetPubKey(keyID)
	retListTags := fakeListTags()
	retDescribeKey := fakeDescribeKey(keyID)

	s.Run("should create a new key successfully", func() {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), alias(id), gomock.Any(), gomock.Any()).Return(&retCreateKey, nil)
		s.mockKmsClient.EXPECT().DescribeKey(ctx, alias(id)).Return(retDescribeKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, keyID).Return(retGetPubKey, nil)
		s.mockKmsClient.EXPECT().ListTags(ctx, keyID, "").Return(retListTags, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.DeletedAt)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.ObjectsAreEqualValues(testutils2.FakeTags(), key.Tags)
	})

	s.Run("should fail with same error if CreateKey fails", func() {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if any function of Get fails", func() {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&retCreateKey, nil)
		s.mockKmsClient.EXPECT().DescribeKey(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *awsKeyStoreTestSuite) TestImport() {
	ctx := context.Background()

	s.Run("should return NotSupportedError", func() {
		_, err := s.keyStore.Import(ctx, "my-id", []byte(""), testutils2.FakeAlgorithm(), testutils2.FakeAttributes())
		assert.Equal(s.T(), errors.ErrNotSupported, err)
	})
}

func (s *awsKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	msg := []byte("some sample message")
	asn1Signature, _ := base64.StdEncoding.DecodeString("MEUCIQDtudqysJc4npK9OCT5whzsE/pZ2zc2DjV9djKcUd1YcwIgHpxvfBLwuQGNu+RbrBq4Skhd9NDQJWo9D2tcsDWRluw=")
	expectedSignature, _ := base64.StdEncoding.DecodeString("7bnasrCXOJ6SvTgk+cIc7BP6Wds3Ng41fXYynFHdWHMenG98EvC5AY275FusGrhKSF300NAlaj0Pa1ywNZGW7A==")
	key := testutils2.FakeKey()
	algo := testutils2.FakeAlgorithm()

	retSign := kms.SignOutput{
		KeyId:     aws.String(keyID),
		Signature: asn1Signature,
	}

	s.Run("should sign a sample message successfully", func() {
		s.mockKeys.EXPECT().Get(ctx, id).Return(key, nil)
		s.mockKmsClient.EXPECT().Sign(ctx, key.Annotations.AWSKeyID, msg, kms.SigningAlgorithmSpecEcdsaSha256).Return(&retSign, nil)

		signature, err := s.keyStore.Sign(ctx, id, msg, algo)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), expectedSignature, signature)
	})

	s.Run("should fail with same error if Get fails", func() {
		s.mockKeys.EXPECT().Get(ctx, id).Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, msg, algo)
		assert.Empty(s.T(), signature)

		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Sign fails", func() {
		s.mockKeys.EXPECT().Get(ctx, id).Return(key, nil)
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), key.Annotations.AWSKeyID, msg, kms.SigningAlgorithmSpecEcdsaSha256).Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, msg, algo)
		assert.Empty(s.T(), signature)

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *awsKeyStoreTestSuite) TestDelete() {
	ctx := context.Background()

	s.Run("should return NotSupportedError", func() {
		err := s.keyStore.Delete(ctx, "my-id")
		assert.Equal(s.T(), errors.ErrNotSupported, err)
	})
}

func (s *awsKeyStoreTestSuite) TestUndelete() {
	ctx := context.Background()

	s.Run("should return NotSupportedError", func() {
		err := s.keyStore.Undelete(ctx, "my-id")
		assert.Equal(s.T(), errors.ErrNotSupported, err)
	})
}

func (s *awsKeyStoreTestSuite) TestDestroy() {
	ctx := context.Background()
	key := testutils2.FakeKey()

	s.Run("should destroy a key successfully", func() {
		s.mockKeys.EXPECT().Get(gomock.Any(), id).Return(key, nil)
		s.mockKmsClient.EXPECT().DeleteKey(gomock.Any(), key.Annotations.AWSKeyID).Return(&kms.ScheduleKeyDeletionOutput{}, nil)

		err := s.keyStore.Destroy(ctx, id)

		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Get fails", func() {
		s.mockKeys.EXPECT().Get(gomock.Any(), id).Return(nil, expectedErr)

		err := s.keyStore.Destroy(ctx, id)

		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if DeleteKey fails", func() {
		s.mockKeys.EXPECT().Get(gomock.Any(), id).Return(key, nil)
		s.mockKmsClient.EXPECT().DeleteKey(gomock.Any(), key.Annotations.AWSKeyID).Return(nil, expectedErr)

		err := s.keyStore.Destroy(ctx, id)

		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *awsKeyStoreTestSuite) TestEncrypt() {
	ctx := context.Background()

	s.Run("should return NotImplementedError", func() {
		_, err := s.keyStore.Encrypt(ctx, "my-id", []byte(""))
		assert.Equal(s.T(), errors.ErrNotImplemented, err)
	})
}

func (s *awsKeyStoreTestSuite) TestDecrypt() {
	ctx := context.Background()

	s.Run("should return NotImplementedError", func() {
		_, err := s.keyStore.Decrypt(ctx, "my-id", []byte(""))
		assert.Equal(s.T(), errors.ErrNotImplemented, err)
	})
}

func ToKmsTags(tags map[string]string) []*kms.Tag {
	var fakeSecretsTags []*kms.Tag

	for key, value := range tags {
		k, v := key, value
		var in = kms.Tag{
			TagKey:   &k,
			TagValue: &v,
		}
		fakeSecretsTags = append(fakeSecretsTags, &in)
	}
	return fakeSecretsTags
}

func fakeDescribeKey(keyID string) *kms.DescribeKeyOutput {
	myArn := "my-key-arn"
	myClusterHsmID := "my-cluster-hsm"
	myAccountID := "my-account"
	myCustomerKeyStoreID := "my-customer-Store"

	return &kms.DescribeKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId:             &keyID,
			Arn:               &myArn,
			Enabled:           aws.Bool(true),
			CreationDate:      aws.Time(time.Now().AddDate(-1, 0, 0)),
			DeletionDate:      aws.Time(time.Now().AddDate(1, 0, 0)),
			ValidTo:           aws.Time(time.Now().AddDate(3, 0, 0)),
			CloudHsmClusterId: &myClusterHsmID,
			CustomKeyStoreId:  &myCustomerKeyStoreID,
			AWSAccountId:      &myAccountID,
		},
	}
}

func fakeListTags() *kms.ListResourceTagsOutput {
	truncatedTagList := false

	return &kms.ListResourceTagsOutput{
		Truncated: &truncatedTagList,
		Tags:      ToKmsTags(testutils2.FakeTags()),
	}
}

func fakeGetPubKey(keyID string) *kms.GetPublicKeyOutput {
	asn1pubKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE1+2SGHa+m/rlkGybiYYuHOvvpINsXmpLJx0OGT5MCrLECKFEKP9VvgWzz2QvoBUWPFhd+jeNb+XbZLxexJr7ZQ==")

	return &kms.GetPublicKeyOutput{
		KeyId:                 &keyID,
		PublicKey:             asn1pubKey,
		KeyUsage:              aws.String(kms.KeyUsageTypeSignVerify),
		CustomerMasterKeySpec: aws.String(kms.CustomerMasterKeySpecEccSecgP256k1),
	}
}
