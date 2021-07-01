package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws/mocks"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	logtestutils "github.com/consensysquorum/quorum-key-manager/pkg/log/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id = "my-key"
)

type awsKeyStoreTestSuite struct {
	suite.Suite
	mockKmsClient *mocks.MockKmsClient
	keyStore      keys.Store
}

func TestAWSKeyStore(t *testing.T) {
	s := new(awsKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *awsKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockKmsClient = mocks.NewMockKmsClient(ctrl)
	s.keyStore = New(s.mockKmsClient, logtestutils.NewMockLogger(ctrl))
}

func (s *awsKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	keyID := "key-id"

	retCreateKey := kms.CreateKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId: aws.String(keyID),
		},
	}

	s.Run("should create a new key successfully", func() {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), alias(id), gomock.Any(), gomock.Any()).Return(&retCreateKey, nil)
		s.getKeyMockCalls(ctx, keyID)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.DeletedAt)
		assert.False(s.T(), key.Metadata.Disabled)
		assert.Equal(s.T(), entities.Ecdsa, key.Algo.Type)
		assert.Equal(s.T(), entities.Secp256k1, key.Algo.EllipticCurve)
		assert.ObjectsAreEqualValues(testutils.FakeTags(), key.Tags)
	})

	s.Run("should fail with same error if CreateKey fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), err, expectedErr)
	})

	s.Run("should fail with same error if any function of Get fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&retCreateKey, nil)
		s.getKeyMockCallsErr(expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), err, expectedErr)
	})
}

func (s *awsKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	keyID := "key_ID"
	expectedPubKey, _ := base64.StdEncoding.DecodeString("BNftkhh2vpv65ZBsm4mGLhzr76SDbF5qSycdDhk+TAqyxAihRCj/Vb4Fs89kL6AVFjxYXfo3jW/l22S8XsSa+2U=")

	retGetPubKey := fakeGetPubKey(keyID)
	retListTags := fakeListTags()
	retDescribeKey := fakeDescribeKey(keyID)

	s.Run("should get a key successfully", func() {
		s.mockKmsClient.EXPECT().DescribeKey(ctx, alias(id)).Return(retDescribeKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, keyID).Return(retGetPubKey, nil)
		s.mockKmsClient.EXPECT().ListTags(ctx, keyID, "").Return(retListTags, nil)

		key, err := s.keyStore.Get(ctx, id)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), key.PublicKey, expectedPubKey)
		assert.ObjectsAreEqualValues(testutils.FakeTags(), key.Tags)
		assert.Equal(s.T(), *retDescribeKey.KeyMetadata.Arn, key.Annotations[awsARN])
		assert.Equal(s.T(), *retDescribeKey.KeyMetadata.AWSAccountId, key.Annotations[awsAccountID])
		assert.Equal(s.T(), *retDescribeKey.KeyMetadata.CustomKeyStoreId, key.Annotations[awsCustomKeyStoreID])
		assert.Equal(s.T(), *retDescribeKey.KeyMetadata.CloudHsmClusterId, key.Annotations[awsCloudHsmClusterID])
	})

	s.Run("should fail with same error if DescribeKey fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().DescribeKey(ctx, gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), err, expectedErr)
	})

	s.Run("should fail with same error if GetPublicKey fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().DescribeKey(ctx, gomock.Any()).Return(retDescribeKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), err, expectedErr)
	})

	s.Run("should fail with same error if ListTags fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().DescribeKey(ctx, gomock.Any()).Return(retDescribeKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, gomock.Any()).Return(retGetPubKey, nil)
		s.mockKmsClient.EXPECT().ListTags(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, id)
		assert.Nil(s.T(), key)

		assert.Equal(s.T(), err, expectedErr)
	})
}

func (s *awsKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	msg := []byte("some sample message")
	myKeyID := "the_ID"

	retSign := kms.SignOutput{
		KeyId:     &myKeyID,
		Signature: []byte("signature"),
	}

	s.Run("should sign a sample message", func() {
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), id, msg, kms.SigningAlgorithmSpecEcdsaSha256).
			Return(&retSign, nil)

		signature, err := s.keyStore.Sign(ctx, id, msg)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)

	})

	s.Run("should fail to sign on error", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), id, msg, kms.SigningAlgorithmSpecEcdsaSha256).
			Return(nil, expectedErr)

		signature, err := s.keyStore.Sign(ctx, id, msg)

		assert.Error(s.T(), err)
		assert.Empty(s.T(), signature)

	})
}

func (s *awsKeyStoreTestSuite) TestList() {
	ctx := context.Background()
	key0, key1 := "my-Key0", "my-key1"
	expected := []string{key0, key1}
	secretsList := []*kms.KeyListEntry{{KeyId: &key0}, {KeyId: &key1}}

	s.Run("should list all keys ids successfully", func() {
		listOutput := &kms.ListKeysOutput{
			Keys: secretsList,
		}

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), 0, "").Return(listOutput, nil)
		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expected, ids)
	})

	s.Run("should list all keys ids successfully with a nextMarker", func() {
		nextMarker := "next"
		listOutput := &kms.ListKeysOutput{
			Keys:       secretsList,
			NextMarker: &nextMarker,
		}

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), 0, nextMarker).Return(listOutput, nil)
		listOutput.NextMarker = nil
		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), 0, "").Return(listOutput, nil)
		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expected, ids)
	})

	s.Run("should return empty keys list if result is nil", func() {
		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), 0, "").Return(&kms.ListKeysOutput{}, nil)
		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Empty(s.T(), ids)
	})

	s.Run("should fail if list fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), 0, "").Return(&kms.ListKeysOutput{}, expectedErr)
		ids, err := s.keyStore.List(ctx)

		assert.Nil(s.T(), ids)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *awsKeyStoreTestSuite) TestDelete() {
	ctx := context.Background()
	myDeletedKeyID := "deleteMe"

	outDeletedKey := &kms.DisableKeyOutput{}

	s.Run("should delete/disable one key successfully", func() {
		s.mockKmsClient.EXPECT().DeleteKey(gomock.Any(), myDeletedKeyID).Return(outDeletedKey, nil)

		err := s.keyStore.Delete(ctx, myDeletedKeyID)

		assert.NoError(s.T(), err)
	})

	s.Run("should fail to delete/disable when error", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().DeleteKey(gomock.Any(), myDeletedKeyID).Return(nil, expectedErr)

		err := s.keyStore.Delete(ctx, myDeletedKeyID)

		assert.Error(s.T(), err)
	})

}

func (s *awsKeyStoreTestSuite) getKeyMockCalls(ctx context.Context, keyID string, ) {
	retGetPubKey := fakeGetPubKey(keyID)
	retListTags := fakeListTags()
	retDescribeKey := fakeDescribeKey(keyID)

	s.mockKmsClient.EXPECT().DescribeKey(ctx, alias(id)).Return(retDescribeKey, nil)
	s.mockKmsClient.EXPECT().GetPublicKey(ctx, keyID).Return(retGetPubKey, nil)
	s.mockKmsClient.EXPECT().ListTags(ctx, keyID, "").Return(retListTags, nil)
}

func (s *awsKeyStoreTestSuite) getKeyMockCallsErr(expectedErr error) {
	s.mockKmsClient.EXPECT().DescribeKey(gomock.Any(), gomock.Any()).Return(nil, expectedErr)
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
	myCustomerKeyStoreID := "my-customer-KeyStore"

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
		Tags:      ToKmsTags(testutils.FakeTags()),
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
