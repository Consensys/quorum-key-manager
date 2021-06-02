package aws

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/keys"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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

	s.T().Run("should create a new key successfully", func(t *testing.T) {
		s.mockKmsClient.EXPECT().CreateKey(gomock.Any(), id, gomock.Any(), gomock.Any()).
			Return(&retCreateKey, nil)
		s.mockKmsClient.EXPECT().GetPublicKey(gomock.Any(), id).
			Return(&retGetPub, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(t, err)
		assert.NotEmpty(t, key.Metadata.CreatedAt)
		assert.NotEmpty(t, key.Metadata.DeletedAt)
		assert.False(t, key.Metadata.Disabled)

	})
}

// TestSign Signature test cases
func (s *awsKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	msg := []byte("some sample message")
	myKeyId := "the_id"

	retSign := kms.SignOutput{
		KeyId:     &myKeyId,
		Signature: []byte("signature"),
	}

	s.T().Run("should sign a sample message", func(t *testing.T) {
		s.mockKmsClient.EXPECT().Sign(gomock.Any(), id, msg).
			Return(&retSign, nil)

		signature, err := s.keyStore.Sign(ctx, id, msg)

		assert.NoError(t, err)
		assert.NotEmpty(t, signature)

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
