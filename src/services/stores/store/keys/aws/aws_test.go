package aws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/keys"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
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
	s.keyStore = New(s.mockKmsClient, log.DefaultLogger())
}

// TestCreate Key creation test cases
func (s *awsKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	creationDate := time.Now()
	deletionDate := creationDate.AddDate(1, 0, 0)

	retCreateKey := kms.CreateKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			CreationDate: &creationDate,
			DeletionDate: &deletionDate,
			Enabled:      aws.Bool(true),
			KeyId:        aws.String("someId"),
		},
	}

	retGetPub := kms.GetPublicKeyOutput{}

	s.Run("should create a new key successfully", func() {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), id, gomock.Any(), gomock.Any()).
			Return(&retCreateKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(gomock.Any(), id).
			Return(&retGetPub, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), key.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), key.Metadata.DeletedAt)
		assert.False(s.T(), key.Metadata.Disabled)

	})

	s.Run("should fail on CreateKey error", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), id, gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), key)

	})

	s.Run("should fail om GetPublicKey error", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), id, gomock.Any(), gomock.Any()).
			Return(&retCreateKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(gomock.Any(), id).
			Return(nil, expectedErr)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), key)

	})
}

// TestGet Get key test cases
func (s *awsKeyStoreTestSuite) TestGet() {
	ctx := context.Background()
	myKeyID := "get_key_ID"
	expectedPubKey := []byte("any value")
	truncatedTagList := false

	retGetPubKey := &kms.GetPublicKeyOutput{KeyId: &myKeyID,
		PublicKey: expectedPubKey}

	retListTags := &kms.ListResourceTagsOutput{
		Truncated: &truncatedTagList,
		Tags:      ToKmsTags(testutils.FakeTags()),
	}

	s.Run("should get a key successfully", func() {
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, myKeyID).Return(retGetPubKey, nil)
		s.mockKmsClient.EXPECT().ListTags(ctx, myKeyID, "").Return(retListTags, nil)

		key, err := s.keyStore.Get(ctx, myKeyID)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), key.PublicKey, expectedPubKey)
		assert.ObjectsAreEqualValues(testutils.FakeTags(), key.Tags)

	})

	s.Run("should fail when listTags failed", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, myKeyID).Return(retGetPubKey, nil)
		s.mockKmsClient.EXPECT().ListTags(ctx, myKeyID, "").Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, myKeyID)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), key)
	})

	s.Run("should fail when getPublicKey failed", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().GetPublicKey(ctx, myKeyID).Return(nil, expectedErr)

		key, err := s.keyStore.Get(ctx, myKeyID)

		assert.Error(s.T(), err)
		assert.Nil(s.T(), key)
	})
}

// TestSign Signature test cases
func (s *awsKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	msg := []byte("some sample message")
	myKeyID := "the_ID"

	retSign := kms.SignOutput{
		KeyId:     &myKeyID,
		Signature: []byte("signature"),
	}

	s.Run("should sign a sample message", func() {
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), id, msg).
			Return(&retSign, nil)

		signature, err := s.keyStore.Sign(ctx, id, msg)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)

	})

	s.Run("should fail to sign on error", func() {
		expectedErr := fmt.Errorf("error")
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), id, msg).
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

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), int64(0), "").Return(listOutput, nil)
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

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), int64(0), "").Return(listOutput, nil)
		listOutput.NextMarker = nil
		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), int64(0), nextMarker).Return(listOutput, nil)
		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expected, ids)
	})

	s.Run("should return empty keys list if result is nil", func() {
		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), int64(0), "").Return(&kms.ListKeysOutput{}, nil)
		ids, err := s.keyStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Empty(s.T(), ids)
	})

	s.Run("should fail if list fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockKmsClient.EXPECT().ListKeys(gomock.Any(), int64(0), "").Return(&kms.ListKeysOutput{}, expectedErr)
		ids, err := s.keyStore.List(ctx)

		assert.Nil(s.T(), ids)
		assert.Equal(s.T(), expectedErr, err)
	})
}

// TestDelete Deletion / Disable key test cases
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
